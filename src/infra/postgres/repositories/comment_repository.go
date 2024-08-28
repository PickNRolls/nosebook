package repos

import (
	"database/sql"
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

type commentLike struct {
	UserId    uuid.UUID `db:"user_id"`
	CommentId uuid.UUID `db:"comment_id"`
}

type commentDest struct {
	id        uuid.UUID    `db:"id"`
	authorId  uuid.UUID    `db:"author_id"`
	postId    uuid.UUID    `db:"post_id"`
	message   string       `db:"message"`
	createdAt time.Time    `db:"created_at"`
	removedAt sql.NullTime `db:"removed_at"`
}

func (this *commentDest) toDomain() *comments.Comment {
	builder := comments.NewBuilder().
		Id(this.id).
		AuthorId(this.authorId).
		Message(this.message).
		PostId(this.postId).
		CreatedAt(this.createdAt)

	if this.removedAt.Valid {
		builder.RemovedAt(this.removedAt.Time)
	}

	return builder.Build()
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

func encodeCursor(comment *comments.Comment) string {
	return fmt.Sprintf(`%v/%v`, comment.Id, comment.CreatedAt.Format(time.RFC3339Nano))
}

func decodeCursor(cursor string) (time.Time, uuid.UUID, *errors.Error) {
	substrings := strings.Split(cursor, "/")
	id := substrings[0]
	createdAt := substrings[1]

	uuidId, err := uuid.Parse(id)
	if err != nil {
		return time.Time{}, uuid.UUID{}, errors.New("DecodeCursorError", "Invalid cursor")
	}

	timestamp, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return time.Time{}, uuid.UUID{}, errors.New("DecodeCursorError", "Invalid cursor")
	}

	return timestamp, uuidId, nil
}

func (repo *CommentRepository) findComments(postIds []uuid.UUID) *generics.BatchQueryResult[*comments.Comment] {
	var result generics.BatchQueryResult[*comments.Comment]
	qb := postgres.NewSquirrel()
	limit := 5

	query := qb.Select(
		append(select_columns, "post_id")...,
	).FromSelect(qb.Select(
		append(
			select_columns,
			"post_id",
			"ROW_NUMBER() OVER (PARTITION BY post_id ORDER BY post_id ASC, created_at ASC, id ASC)",
		)...,
	).From(comments_table).Join(
		post_comments_table+" on c.id = pc.comment_id",
	).Where(
		"removed_at IS NULL",
	).Where(squirrel.Eq{"post_id": postIds}).OrderBy(
		"post_id ASC, created_at ASC, id ASC",
	).Limit(uint64(limit+1)), "subquery").Where(
		"row_number <= ?", limit+1,
	)

	var commentsRows []*comments.Comment
	sql, args, _ := query.ToSql()

	err := repo.db.Select(&commentsRows, sql, args...)
	if err != nil {
		result.Err = errors.New("FindError", err.Error())
		return &result
	}

	for _, row := range commentsRows {
		if !result.HasEntry(row.PostId) {
			result.AddEntryOnce(row.PostId)
		}

		single := result.SingleResultOf(row.PostId)
		if len(single.Data) < limit {
			single.Data = append(single.Data, row)
		} else {
			single.Next = encodeCursor(single.Data[len(single.Data)-1])
		}
	}

	return &result
}

func (repo *CommentRepository) findCommentLikes(coms []*comments.Comment) ([]*commentLike, *errors.Error) {
	qb := postgres.NewSquirrel()
	commentIds := make([]uuid.UUID, len(coms))
	for _, comment := range coms {
		commentIds = append(commentIds, comment.Id)
	}

	if len(commentIds) > 0 {
		var likes []*commentLike

		sql, args, _ := qb.Select("user_id", "comment_id").From("comment_likes").Where(squirrel.Eq{"comment_id": commentIds}).ToSql()
		err := repo.db.Select(&likes, sql, args...)
		if err != nil {
			return nil, errors.New("FindError", err.Error())
		}

		return likes, nil
	}

	return make([]*commentLike, 0), nil
}

func (repo *CommentRepository) FindByPostIds(postIds []uuid.UUID) *generics.BatchQueryResult[*comments.Comment] {
	return repo.findComments(postIds)
}

func (repo *CommentRepository) FindByFilter(filter structs.QueryFilter, limitPointer *uint) *generics.SingleQueryResult[*comments.Comment] {
	qb := postgres.NewSquirrel()
	limit := uint(20)
	if limitPointer != nil {
		limit = *limitPointer
	}

	var result generics.SingleQueryResult[*comments.Comment]

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
		timestamp, id, err := decodeCursor(filter.Next)
		if err != nil {
			result.Err = err
			return &result
		}

		cursorQuery = cursorQuery.Where("(created_at, id) > (?, ?)", timestamp, id)
	}

	if filter.Prev != "" {
		timestamp, id, err := decodeCursor(filter.Prev)
		if err != nil {
			result.Err = err
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
			result.Prev = encodeCursor(firstComment)
		}

		if nextCount.Count > 0 {
			result.Next = encodeCursor(lastComment)
		}
	}

	// likes, error := repo.findCommentLikes(result.Data)
	// if error != nil {
	// 	result.Err = error
	// 	return &result
	// }
	//
	// for _, like := range likes {
	// 	for _, comment := range result.Data {
	// 		if like.CommentId == comment.Id {
	// 			comment.LikedBy = append(comment.LikedBy, like.UserId)
	// 		}
	// 	}
	// }

	return &result
}

func (repo *CommentRepository) FindById(id uuid.UUID, includeRemoved bool) *comments.Comment {
	qb := postgres.NewSquirrel()
	dest := commentDest{}

	{
		query := qb.Select(append(select_columns, "post_id")...).
			From(comments_table).
			Join(
				post_comments_table+" on c.id = pc.comment_id",
			).
			Where("id = ?", id)

		if !includeRemoved {
			query = query.Where("removed_at IS NULL")
		}
		sql, args, _ := query.ToSql()
		err := repo.db.Get(&dest, sql, args...)

		if err != nil {
			return nil
		}
	}

	return dest.toDomain()
}

func (repo *CommentRepository) Save(comment *comments.Comment) *errors.Error {
	qb := postgres.NewSquirrel()
	tx, err := repo.db.Beginx()
	if err != nil {
		return errors.From(err)
	}

	for _, event := range comment.Events() {
		switch event.Type() {
		case comments.CREATED:
			{
				sql, args, _ := qb.Insert(comments_table).Columns(
					insert_columns...,
				).Values(
					comment.Id, comment.AuthorId, comment.Message, comment.CreatedAt, comment.RemovedAt,
				).ToSql()
				_, err = tx.Exec(sql, args...)

				if err != nil {
					tx.Rollback()
					return errors.From(err)
				}
			}

			if comment.PostId != uuid.Nil {
				sql, args, _ := qb.Insert("post_comments").Columns("post_id", "comment_id").Values(comment.PostId, comment.Id).ToSql()
				_, err = tx.Exec(sql, args...)

				if err != nil {
					tx.Rollback()
					return errors.From(err)
				}
			}

		case comments.REMOVED:
			sql, args, _ := qb.Update(comments_table).Set("removed_at", comment.RemovedAt).Where("id = ?", comment.Id).ToSql()
			_, err := tx.Exec(sql, args...)
			if err != nil {
				tx.Rollback()
				return errors.From(err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.From(err)
	}

	return nil
}
