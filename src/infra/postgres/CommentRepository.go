package postgres

import (
	"nosebook/src/domain/comments"
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

func (repo *CommentRepository) CreateForPost(postId uuid.UUID, comment *comments.Comment) (*comments.Comment, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(`INSERT INTO comments (
    id,
    author_id,
    message,
    created_at
  ) VALUES (
    $1,
    $2,
    $3,
    $4
  )`, comment.Id, comment.AuthorId, comment.Message, comment.CreatedAt)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(`INSERT INTO post_comments (
    post_id,
    comment_id
	) VALUES (
    $1,
    $2
	)`, postId, comment.Id)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return comment, nil
}
