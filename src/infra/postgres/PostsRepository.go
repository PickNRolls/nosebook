package postgres

import (
	"nosebook/src/domain/posts"
	"nosebook/src/services/posting/interfaces"

	"github.com/google/uuid"
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

func (repo *PostsRepository) FindById(id uuid.UUID) *posts.Post {
	post := posts.Post{}
	err := repo.db.Get(&post, `SELECT
		id,
		author_id,
		owner_id,
		message,
		created_at
			FROM posts WHERE
		id = $1
	`, id)

	if err != nil {
		return nil
	}

	return &post
}

func (repo *PostsRepository) Create(post *posts.Post) (*posts.Post, error) {
	_, err := repo.db.NamedExec(`INSERT INTO posts (
	  id,
	  author_id,
	  owner_id,
	  message,
	  created_at
	) VALUES (
	  :id,
	  :author_id,
	  :owner_id,
	  :message,
	  :created_at
	)`, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (repo *PostsRepository) Remove(post *posts.Post) (*posts.Post, error) {
	_, err := repo.db.NamedExec(`DELETE FROM posts WHERE
		id = :id	
	`, post)
	if err != nil {
		return nil, err
	}

	return post, nil
}
