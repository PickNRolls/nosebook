package repos

import (
	"nosebook/src/domain/comments"
	"nosebook/src/infra/postgres"
	"nosebook/src/services/commenting/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) interfaces.CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (repo *CommentRepository) FindById(id uuid.UUID, includeRemoved bool) *comments.Comment {
	qb := postgres.NewSquirrel()
	comment := comments.Comment{}

	{
		query := qb.Select("id", "author_id", "message", "created_at", "removed_at").From("comments").Where("id = (?)", id)
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
			sql, args, _ := qb.Update("comments").Set("removed_at", comment.RemovedAt).Where("id = ?", comment.Id).ToSql()
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
		sql, args, _ := qb.Insert("comments").Columns(
			"id", "author_id", "message", "created_at", "removed_at",
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
