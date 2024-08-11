package repos

import (
	"fmt"
	"nosebook/src/domain/comments"
	"nosebook/src/errors"
	"nosebook/src/generics"
	"nosebook/src/infra/postgres"
	"nosebook/src/services/commenting/interfaces"
	"nosebook/src/services/commenting/structs"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

var select_columns = []string{"id", "author_id", "message", "created_at", "removed_at"}
var insert_columns = select_columns

const comments_table = "comments as c"
const post_comments_table = "post_comments as pc"

func NewCommentRepository(db *sqlx.DB) interfaces.CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (repo *CommentRepository) FindByFilter(filter structs.QueryFilter, limitPointer *uint) *generics.QuerySingleResult[*comments.Comment] {
	qb := postgres.NewSquirrel()
	limit := uint(20)
	if limitPointer != nil {
		limit = *limitPointer
	}

	var result generics.QuerySingleResult[*comments.Comment]

	if filter.PostId == uuid.Nil && filter.AuthorId == uuid.Nil {
		result.Err = errors.New("FindError", "You must specify either postId or authorId at least")
		return &result
	}

	if filter.Next != "" && filter.Prev != "" {
		result.Err = errors.New("FindError", "You can't specify next and prev at the same time")
		return &result
	}

	if filter.Last {
		if filter.Next != "" || filter.Prev != "" {
			result.Err = errors.New("FindError", "You can't specify next/prev and last at the same time")
			return &result
		}
	}

	var commentsRows []*comments.Comment
	query := qb.Select(select_columns...).From(comments_table).Where("removed_at IS NULL")

	if filter.PostId != uuid.Nil {
		query = query.Columns("post_id").Join(
			post_comments_table+" on c.id = pc.comment_id",
		).Where("post_id = ?", filter.PostId)
	}

	if filter.AuthorId != uuid.Nil {
		query = query.Where("author_id = ?", filter.AuthorId)
	}

	limitQuery := query.Limit(uint64(limit))
	lastQuery := qb.Select("*").FromSelect(
		limitQuery.OrderBy("created_at DESC, id ASC"), "last",
	).OrderBy("created_at ASC, id ASC")
	cursorQuery := limitQuery.OrderBy("created_at ASC, id ASC")

	if filter.Next != "" {
		substrings := strings.Split(filter.Next, "/")
		id := substrings[0]
		createdAt := substrings[1]

		timestamp, err := time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			result.Err = errors.New("FindError", "Invalid cursor")
			return &result
		}

		cursorQuery = cursorQuery.Where("(created_at, id) > (?, ?)", timestamp, id)
	}

	if filter.Prev != "" {
		substrings := strings.Split(filter.Prev, "/")
		id := substrings[0]
		createdAt := substrings[1]

		timestamp, err := time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			result.Err = errors.New("FindError", "Invalid cursor")
			return &result
		}

		cursorQuery = cursorQuery.Where("(created_at, id) < (?, ?)", timestamp, id)
	}

	resultQuery := cursorQuery
	if filter.Last {
		resultQuery = lastQuery
	}

	sql, args, _ := resultQuery.ToSql()
	err := repo.db.Select(&commentsRows, sql, args...)

	if err != nil {
		result.Err = errors.New("FindError", err.Error())
		return &result
	}

	result.Data = commentsRows

	if len(commentsRows) > 0 {
		var nextCount struct {
			Count int `db:"count"`
		}

		var prevCount struct {
			Count int `db:"count"`
		}

		firstComment := commentsRows[0]
		lastComment := commentsRows[len(commentsRows)-1]
		countQuery := query.RemoveColumns().Columns("count(*)")

		countPrevQuery := countQuery.Where("(created_at, id) < (?, ?)", firstComment.CreatedAt, firstComment.Id)
		sql, args, _ := countPrevQuery.ToSql()
		err := repo.db.Get(&prevCount, sql, args...)
		if err != nil {
			result.Err = errors.New("FindError", err.Error())
			return &result
		}

		countNextQuery := countQuery.Where("(created_at, id) > (?, ?)", lastComment.CreatedAt, lastComment.Id)
		sql, args, _ = countNextQuery.ToSql()
		err = repo.db.Get(&nextCount, sql, args...)
		if err != nil {
			result.Err = errors.New("FindError", err.Error())
			return &result
		}

		if prevCount.Count > 0 {
			result.Prev = fmt.Sprintf(`%v/%v`, firstComment.Id, firstComment.CreatedAt.Format(time.RFC3339Nano))
		}

		if nextCount.Count > 0 {
			result.Next = fmt.Sprintf(`%v/%v`, lastComment.Id, lastComment.CreatedAt.Format(time.RFC3339Nano))
		}
	}

	commentIds := make([]uuid.UUID, 0)
	for _, comment := range result.Data {
		commentIds = append(commentIds, comment.Id)
	}

	if len(commentIds) > 0 {
		var likes []struct {
			UserId    uuid.UUID `db:"user_id"`
			CommentId uuid.UUID `db:"comment_id"`
		}

		sql, args, _ := qb.Select("user_id", "comment_id").From("comment_likes").Where(squirrel.Eq{"comment_id": commentIds}).ToSql()
		err = repo.db.Select(&likes, sql, args...)
		if err != nil {
			result.Err = errors.New("FindError", err.Error())
			return &result
		}

		for _, like := range likes {
			for _, comment := range result.Data {
				if like.CommentId == comment.Id {
					comment.LikedBy = append(comment.LikedBy, like.UserId)
				}
			}
		}
	}

	return &result
}

func (repo *CommentRepository) FindById(id uuid.UUID, includeRemoved bool) *comments.Comment {
	qb := postgres.NewSquirrel()
	comment := comments.Comment{}

	{
		query := qb.Select(select_columns...).From(comments_table).Where("id = (?)", id)
		if !includeRemoved {
			query = query.Where("removed_at IS NULL")
		}
		sql, args, _ := query.ToSql()
		err := repo.db.Get(&comment, sql, args...)

		if err != nil {
			return nil
		}
	}

	{
		var likes []struct {
			UserId uuid.UUID `db:"user_id"`
		}

		sql, args, _ := qb.Select("user_id").From("comment_likes").Where("comment_id = (?)", id).ToSql()
		err := repo.db.Select(&likes, sql, args...)

		if err != nil {
			return nil
		}

		for _, like := range likes {
			comment.LikedBy = append(comment.LikedBy, like.UserId)
		}
	}

	return &comment
}

func (repo *CommentRepository) Save(comment *comments.Comment) (*comments.Comment, error) {
	qb := postgres.NewSquirrel()
	tx, err := repo.db.Beginx()
	if err != nil {
		return nil, err
	}

	for _, event := range comment.Events {
		switch event.Type() {
		case comments.LIKED:
			likeEvent := event.(*comments.CommentLikeEvent)
			sql, args, _ := qb.Insert("comment_likes").Columns("user_id", "comment_id").Values(likeEvent.UserId, comment.Id).ToSql()
			_, err := tx.Exec(sql, args...)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

		case comments.UNLIKED:
			unlikeEvent := event.(*comments.CommentUnlikeEvent)
			sql, args, _ := qb.Delete("comment_likes").Where("user_id = ? AND comment_id = ?", unlikeEvent.UserId, comment.Id).ToSql()
			_, err := tx.Exec(sql, args...)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

		case comments.REMOVED:
			sql, args, _ := qb.Update(comments_table).Set("removed_at", comment.RemovedAt).Where("id = ?", comment.Id).ToSql()
			_, err := tx.Exec(sql, args...)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (repo *CommentRepository) Create(comment *comments.Comment) (*comments.Comment, error) {
	qb := postgres.NewSquirrel()
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}

	{
		sql, args, _ := qb.Insert(comments_table).Columns(
			insert_columns...,
		).Values(
			comment.Id, comment.AuthorId, comment.Message, comment.CreatedAt, comment.RemovedAt,
		).ToSql()
		_, err = tx.Exec(sql, args...)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if comment.PostId != uuid.Nil {
		sql, args, _ := qb.Insert("post_comments").Columns("post_id", "comment_id").Values(comment.PostId, comment.Id).ToSql()
		_, err = tx.Exec(sql, args...)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return comment, nil
}
