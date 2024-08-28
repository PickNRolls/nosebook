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
	builder := posts.NewBuilder().
		Id(post.id).
		AuthorId(post.authorId).
		OwnerId(post.ownerId).
		Message(post.message).
		CreatedAt(post.createdAt)

	if post.removedAt.Valid {
		builder.RemovedAt(post.removedAt.Time)
	}

	return builder.Build()
}

func NewPostsRepository(db *sqlx.DB) interfaces.PostRepository {
	return &PostsRepository{
		db: db,
	}
}

var posts_table = "posts"
var posts_select_columns = []string{"id", "author_id", "owner_id", "message", "created_at", "removed_at"}
var posts_insert_columns = posts_select_columns

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
