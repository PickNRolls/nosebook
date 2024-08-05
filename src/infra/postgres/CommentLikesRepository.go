package postgres

import (
	"nosebook/src/domain/comments"
	"nosebook/src/services/commenting/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentLikesRepository struct {
	db *sqlx.DB
}

func NewCommentLikesRepository(db *sqlx.DB) interfaces.CommentLikesRepository {
	return &CommentLikesRepository{
		db: db,
	}
}

func (repo *CommentLikesRepository) Find(commentId uuid.UUID, userId uuid.UUID) *comments.CommentLike {
	like := comments.CommentLike{}
	err := repo.db.Get(&like, `SELECT
		author_id,
		comment_id
	FROM
		comment_likes
	WHERE
		author_id = $1 AND comment_id = $2
	`, userId, commentId)

	if err != nil {
		return nil
	}

	return &like
}

func (repo *CommentLikesRepository) Create(like *comments.CommentLike) (*comments.CommentLike, error) {
	_, err := repo.db.NamedExec(`INSERT INTO comment_likes (
	  author_id,
	  comment_id
	) VALUES (
	  :author_id,
	  :comment_id
	)`, like)
	if err != nil {
		return nil, err
	}

	return like, nil
}

func (repo *CommentLikesRepository) Remove(like *comments.CommentLike) (*comments.CommentLike, error) {
	_, err := repo.db.NamedExec(`DELETE FROM comment_likes WHERE
		author_id = :author_id AND
		comment_id = :comment_id
	`, like)
	if err != nil {
		return nil, err
	}

	return like, nil
}
