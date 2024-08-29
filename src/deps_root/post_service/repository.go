package rootpostservice

import (
	"database/sql"
	"nosebook/src/domain/post"
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var posts_table = "posts"
var posts_select_columns = []string{"id", "author_id", "owner_id", "message", "created_at", "removed_at"}
var posts_insert_columns = posts_select_columns

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (this *Repository) FindById(id uuid.UUID) *domainpost.Post {
	qb := postgres.NewSquirrel()

	dest := struct {
		Id        uuid.UUID    `db:"id"`
		AuthorId  uuid.UUID    `db:"author_id"`
		OwnerId   uuid.UUID    `db:"owner_id"`
		Message   string       `db:"message"`
		CreatedAt time.Time    `db:"created_at"`
		RemovedAt sql.NullTime `db:"removed_at"`
	}{}

	sqlQuery, args, _ := qb.Select(
		posts_select_columns...,
	).From(
		posts_table,
	).Where(
		"id = ? AND removed_at IS NULL",
		id,
	).ToSql()

	err := this.db.Get(&dest, sqlQuery, args...)
	if err != nil {
		return nil
	}

	builder := domainpost.NewBuilder().
		Id(dest.Id).
		AuthorId(dest.AuthorId).
		OwnerId(dest.OwnerId).
		Message(dest.Message).
		CreatedAt(dest.CreatedAt)

	if dest.RemovedAt.Valid {
		builder.RemovedAt(dest.RemovedAt.Time)
	}

	return builder.Build()
}

func (this *Repository) Save(post *domainpost.Post) *errors.Error {
	qb := postgres.NewSquirrel()

	for _, event := range post.Events() {
		switch event.Type() {
		case domainpost.CREATED:
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
			_, err := this.db.Exec(sql, args...)

			if err != nil {
				return errors.From(err)
			}

		case domainpost.EDITED:
			editedEvent := event.(*domainpost.PostEditedEvent)
			sql, args, _ := qb.Update(posts_table).Set(
				"message", editedEvent.Message,
			).Where(
				"id = ?",
				post.Id,
			).ToSql()
			_, err := this.db.Exec(sql, args...)

			if err != nil {
				return errors.From(err)
			}

		case domainpost.REMOVED:
			sql, args, _ := qb.Update(posts_table).Set(
				"removed_at", post.RemovedAt,
			).Where(
				"id = ?",
				post.Id,
			).ToSql()
			_, err := this.db.Exec(sql, args...)

			if err != nil {
				return errors.From(err)
			}
		}
	}

	return nil
}
