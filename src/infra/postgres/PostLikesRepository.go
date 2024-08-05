package postgres

import (
	"nosebook/src/domain/posts"
	"nosebook/src/services/posting/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostLikesRepository struct {
	db *sqlx.DB
}

func NewPostLikesRepository(db *sqlx.DB) interfaces.PostLikesRepository {
	return &PostLikesRepository{
		db: db,
	}
}

func (repo *PostLikesRepository) Find(postId uuid.UUID, userId uuid.UUID) *posts.PostLike {
	like := posts.PostLike{}
	err := repo.db.Get(&like, `SELECT
		author_id,
		post_id
	FROM
		post_likes
	WHERE
		author_id = $1 AND post_id = $2
	`, userId, postId)

	if err != nil {
		return nil
	}

	return &like
}

func (repo *PostLikesRepository) Create(like *posts.PostLike) (*posts.PostLike, error) {
	_, err := repo.db.NamedExec(`INSERT INTO post_likes (
	  author_id,
	  post_id
	) VALUES (
	  :author_id,
	  :post_id
	)`, like)
	if err != nil {
		return nil, err
	}

	return like, nil
}

func (repo *PostLikesRepository) Remove(like *posts.PostLike) (*posts.PostLike, error) {
	_, err := repo.db.NamedExec(`DELETE FROM post_likes WHERE
		author_id = :author_id AND
		post_id = :post_id
	`, like)
	if err != nil {
		return nil, err
	}

	return like, nil
}
