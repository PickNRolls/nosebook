package repos

import (
	"errors"
	"fmt"
	"nosebook/src/domain/posts"
	"nosebook/src/services/posting/interfaces"
	"nosebook/src/services/posting/structs"
	"strings"
	"time"

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

	if filter.Cursor != "" {
		substrings := strings.Split(filter.Cursor, "/")
		id := substrings[0]
		createdAt := substrings[1]

		timestamp, err := time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			result.Err = errors.New("Invalid cursor.")
			return result
		}

		whereClause = append(whereClause, fmt.Sprintf("(created_at, id) < ($%v, $%v)", len(args)+1, len(args)+2))
		args = append(args, timestamp, id)
	}

	where := strings.Join(whereClause, " AND ")
	limit := 10

	var posts []*posts.Post
	err := repo.db.Select(&posts, fmt.Sprintf(`SELECT
		id,
		author_id,
		owner_id,
		message,
		created_at
			FROM posts WHERE
		removed_at IS NULL AND %v
			ORDER BY created_at DESC, id DESC
			LIMIT %v
	`, where, limit), args...)

	if err != nil {
		result.Err = err
		return result
	}

	if filter.Cursor != "" {
		whereClause = whereClause[:len(whereClause)-2]
		args = args[:len(args)-2]
	}

	result.Data = posts

	if len(posts) == limit {
		var remainingCount struct {
			Count int `db:"count"`
		}
		lastPost := posts[len(posts)-1]
		whereWithNextCursor := fmt.Sprintf(`%v AND (created_at, id) < ($%v, $%v)`, where, len(args)+1, len(args)+2)
		args = append(args, lastPost.CreatedAt)
		args = append(args, lastPost.Id)
		err = repo.db.Get(&remainingCount, fmt.Sprintf(`SELECT
			COUNT(*)
				FROM posts WHERE
			removed_at IS NULL AND %v
		`, whereWithNextCursor), args...)

		if err != nil {
			result.Err = err
			return result
		}

		result.RemainingCount = remainingCount.Count
		result.Next = fmt.Sprintf(`%v/%v`, lastPost.Id, lastPost.CreatedAt.Format(time.RFC3339Nano))
	}

	postIds := make([]uuid.UUID, 0)
	for _, post := range result.Data {
		postIds = append(postIds, post.Id)
	}

	if len(postIds) > 0 {
		var likes []struct {
			UserId uuid.UUID `db:"user_id"`
			PostId uuid.UUID `db:"post_id"`
		}

		query, args, err := sqlx.In(`SELECT
		user_id,
		post_id
		  FROM post_likes WHERE
		post_id IN (?)
	`, postIds)
		if err != nil {
			result.Err = err
			return result
		}

		query = repo.db.Rebind(query)
		err = repo.db.Select(&likes, query, args...)
		if err != nil {
			result.Err = err
			return result
		}

		for _, like := range likes {
			for _, post := range result.Data {
				if like.PostId == post.Id {
					post.LikedBy = append(post.LikedBy, like.UserId)
				}
			}
		}
	}

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
		id = $1 AND removed_at IS NULL
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
	_, err := repo.db.NamedExec(`UPDATE posts SET
		removed_at = NOW()
			WHERE
		id = :id	
	`, post)
	if err != nil {
		return nil, err
	}

	post.RemovedAt = time.Now()

	return post, nil
}
