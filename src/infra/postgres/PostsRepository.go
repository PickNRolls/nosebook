package postgres

import (
	"errors"
	"fmt"
	"nosebook/src/domain/posts"
	"nosebook/src/services/posting/interfaces"
	"nosebook/src/services/posting/structs"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostsRepository struct {
	db *sqlx.DB
}

func NewPostsRepository(db *sqlx.DB) interfaces.PostRepository {
	return &PostsRepository{
		db: db,
	}
}

func (repo *PostsRepository) FindByFilter(filter structs.QueryFilter) structs.QueryResult {
	var result structs.QueryResult
	whereClause := make([]string, 0)
	args := make([]any, 0)

	if filter.OwnerId != uuid.Nil {
		whereClause = append(whereClause, fmt.Sprintf("owner_id = $%v", len(args)+1))
		args = append(args, filter.OwnerId)
	}

	if filter.AuthorId != uuid.Nil {
		whereClause = append(whereClause, fmt.Sprintf("author_id = $%v", len(args)+1))
		args = append(args, filter.AuthorId)
	}

	if filter.OwnerId == uuid.Nil && filter.AuthorId == uuid.Nil {
		result.Err = errors.New("You must specify either ownerId or authorId at least.")
		return result
	}

	where := strings.Join(whereClause, " AND ")

	var posts []*posts.Post
	err := repo.db.Select(&posts, fmt.Sprintf(`SELECT
		id,
		author_id,
		owner_id,
		message,
		created_at
			FROM posts WHERE
		%v
			ORDER BY created_at DESC
	`, where), args...)

	if err != nil {
		result.Err = err
		return result
	}

	result.Data = posts

	return result
}

func (repo *PostsRepository) Save(post *posts.Post) (*posts.Post, error) {
	tx, err := repo.db.Beginx()
	if err != nil {
		return nil, err
	}

	for _, event := range post.Events {
		switch event.Type() {
		case posts.LIKED:
			likeEvent := event.(*posts.PostLikeEvent)
			_, err := tx.Exec(`INSERT INTO post_likes (
				user_id,
				post_id
			) VALUES (
				$1,
				$2
			)`, likeEvent.UserId, post.Id)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

		case posts.UNLIKED:
			unlikeEvent := event.(*posts.PostUnlikeEvent)
			_, err := tx.Exec(`DELETE FROM post_likes WHERE
				user_id = $1 AND post_id = $2
			`, unlikeEvent.UserId, post.Id)
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

	return post, nil
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

	var likes []struct {
		UserId uuid.UUID `db:"user_id"`
	}
	err = repo.db.Select(&likes, `SELECT
		user_id
			FROM post_likes WHERE
		post_id = $1
	`, id)

	if err != nil {
		return nil
	}

	for _, like := range likes {
		post.LikedBy = append(post.LikedBy, like.UserId)
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
