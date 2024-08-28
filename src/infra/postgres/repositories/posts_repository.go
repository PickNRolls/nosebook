package repos

import (
	"database/sql"
	"nosebook/src/domain/posts"
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	"nosebook/src/services/posting/interfaces"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostsRepository struct {
	db *sqlx.DB
}

type PostDest struct {
	id        uuid.UUID    `db:"id"`
	authorId  uuid.UUID    `db:"author_id"`
	ownerId   uuid.UUID    `db:"owner_id"`
	message   string       `db:"message"`
	createdAt time.Time    `db:"created_at"`
	removedAt sql.NullTime `db:"removed_at"`
}

func (post *PostDest) ToDomain() *posts.Post {
	return posts.NewPost(
		post.id,
		post.authorId,
		post.ownerId,
		post.message,
		post.createdAt,
		post.removedAt,
		false,
	)
}

func NewPostsRepository(db *sqlx.DB) interfaces.PostRepository {
	return &PostsRepository{
		db: db,
	}
}

var posts_table = "posts"
var posts_select_columns = []string{"id", "author_id", "owner_id", "message", "created_at", "removed_at"}
var posts_insert_columns = posts_select_columns

//	func (repo *PostsRepository) FindByFilter(filter structs.QueryFilter) *generics.SingleQueryResult[*posts.Post] {
//		var result generics.SingleQueryResult[*posts.Post]
//		whereClause := make([]string, 0)
//		args := make([]any, 0)
//
//		if filter.OwnerId != uuid.Nil {
//			whereClause = append(whereClause, fmt.Sprintf("owner_id = $%v", len(args)+1))
//			args = append(args, filter.OwnerId)
//		}
//
//		if filter.AuthorId != uuid.Nil {
//			whereClause = append(whereClause, fmt.Sprintf("author_id = $%v", len(args)+1))
//			args = append(args, filter.AuthorId)
//		}
//
//		if filter.OwnerId == uuid.Nil && filter.AuthorId == uuid.Nil {
//			result.Err = errors.New("FindError", "You must specify either ownerId or authorId at least")
//			return &result
//		}
//
//		if filter.Cursor != "" {
//			substrings := strings.Split(filter.Cursor, "/")
//			id := substrings[0]
//			createdAt := substrings[1]
//
//			timestamp, err := time.Parse(time.RFC3339Nano, createdAt)
//			if err != nil {
//				result.Err = errors.New("FindError", "Invalid cursor")
//				return &result
//			}
//
//			whereClause = append(whereClause, fmt.Sprintf("(created_at, id) < ($%v, $%v)", len(args)+1, len(args)+2))
//			args = append(args, timestamp, id)
//		}
//
//		where := strings.Join(whereClause, " AND ")
//		limit := 10
//
//		var posts []*posts.Post
//		err := repo.db.Select(&posts, fmt.Sprintf(`SELECT
//			id,
//			author_id,
//			owner_id,
//			message,
//			created_at
//				FROM posts WHERE
//			removed_at IS NULL AND %v
//				ORDER BY created_at DESC, id DESC
//				LIMIT %v
//		`, where, limit), args...)
//
//		if err != nil {
//			result.Err = errors.New("FindError", err.Error())
//			return &result
//		}
//
//		if filter.Cursor != "" {
//			whereClause = whereClause[:len(whereClause)-2]
//			args = args[:len(args)-2]
//		}
//
//		result.Data = posts
//
//		if len(posts) == limit {
//			var remainingCount struct {
//				Count int `db:"count"`
//			}
//			lastPost := posts[len(posts)-1]
//			whereWithNextCursor := fmt.Sprintf(`%v AND (created_at, id) < ($%v, $%v)`, where, len(args)+1, len(args)+2)
//			args = append(args, lastPost.CreatedAt)
//			args = append(args, lastPost.Id)
//			err = repo.db.Get(&remainingCount, fmt.Sprintf(`SELECT
//				COUNT(*)
//					FROM posts WHERE
//				removed_at IS NULL AND %v
//			`, whereWithNextCursor), args...)
//
//			if err != nil {
//				result.Err = errors.New("FindError", err.Error())
//				return &result
//			}
//
//			result.RemainingCount = remainingCount.Count
//			result.Next = fmt.Sprintf(`%v/%v`, lastPost.Id, lastPost.CreatedAt.Format(time.RFC3339Nano))
//		}
//
//		postIds := make([]uuid.UUID, 0)
//		for _, post := range result.Data {
//			postIds = append(postIds, post.Id)
//		}
//
//		if len(postIds) > 0 {
//			var likes []struct {
//				UserId uuid.UUID `db:"user_id"`
//				PostId uuid.UUID `db:"post_id"`
//			}
//
//			query, args, err := sqlx.In(`SELECT
//				user_id,
//				post_id
//			  	FROM post_likes WHERE
//				post_id IN (?)
//			`, postIds)
//			if err != nil {
//				result.Err = errors.New("FindError", err.Error())
//				return &result
//			}
//
//			query = repo.db.Rebind(query)
//			err = repo.db.Select(&likes, query, args...)
//			if err != nil {
//				result.Err = errors.New("FindError", err.Error())
//				return &result
//			}
//
//			for _, like := range likes {
//				for _, post := range result.Data {
//					if like.PostId == post.Id {
//						post.LikedBy = append(post.LikedBy, like.UserId)
//					}
//				}
//			}
//		}
//
//		return &result
//	}
func (repo *PostsRepository) Save(post *posts.Post) *errors.Error {
	qb := postgres.NewSquirrel()

	for _, event := range post.Events() {
		switch event.Type() {
		case posts.CREATED:
			sql, args, _ := qb.Insert(posts_table).Columns(
				posts_insert_columns...,
			).Values(
				post.Id,
				post.AuthorId,
				post.OwnerId,
				post.Message,
				post.CreatedAt,
				post.RemovedAt,
			).ToSql()
			_, err := repo.db.Exec(sql, args...)

			if err != nil {
				return errors.From(err)
			}

		case posts.EDITED:
			editedEvent := event.(*posts.PostEditedEvent)
			sql, args, _ := qb.Update(posts_table).Set(
				"message", editedEvent.Message,
			).Where(
				"id = ?",
				post.Id,
			).ToSql()
			_, err := repo.db.Exec(sql, args...)

			if err != nil {
				return errors.From(err)
			}

		case posts.REMOVED:
			sql, args, _ := qb.Update(posts_table).Set(
				"removed_at", post.RemovedAt,
			).Where(
				"id = ?",
				post.Id,
			).ToSql()
			_, err := repo.db.Exec(sql, args...)

			if err != nil {
				return errors.From(err)
			}
		}
	}

	return nil
}

func (repo *PostsRepository) FindById(id uuid.UUID) *posts.Post {
	qb := postgres.NewSquirrel()

	postDest := PostDest{}
	sql, args, _ := qb.Select(
		posts_select_columns...,
	).From(
		posts_table,
	).Where(
		"id = ? AND removed_at IS NULL",
		id,
	).ToSql()

	err := repo.db.Get(&postDest, sql, args...)
	if err != nil {
		return nil
	}

	return postDest.ToDomain()
}
