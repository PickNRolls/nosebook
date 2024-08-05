package postgres

import (
	"nosebook/src/domain/posts"
	"nosebook/src/services/posting/interfaces"

	"github.com/jmoiron/sqlx"
)

type PostsRepository struct {
	db *sqlx.DB
}

func NewPostsRepository(db *sqlx.DB) interfaces.PostsRepository {
	return &PostsRepository{
		db: db,
	}
}

func (repo *PostsRepository) Create(post *posts.Post) (*posts.Post, error) {
	_, err := repo.db.NamedExec(`INSERT INTO posts (
	  id,
	  author_id,
	  owner_id,
	  message
	) VALUES (
	  :id,
	  :author_id,
	  :owner_id,
	  :message
	)`, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}
