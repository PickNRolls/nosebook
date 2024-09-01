package rootcommentservice

import (
	"database/sql"
	"nosebook/src/domain/comment"
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type commentRepository struct {
	db          *sqlx.DB
	permissions domaincomment.Permissions
}

func newCommentRepository(db *sqlx.DB, permissions domaincomment.Permissions) *commentRepository {
	return &commentRepository{
		db:          db,
		permissions: permissions,
	}
}

const comments_table = "comments as c"
const post_comments_table = "post_comments as pc"

var select_columns = []string{"id", "author_id", "message", "created_at", "removed_at"}
var insert_columns = select_columns

func (this *commentRepository) FindById(id uuid.UUID, includeRemoved bool) *domaincomment.Comment {
	qb := postgres.NewSquirrel()

	query := qb.Select(append(select_columns, "post_id")...).
		From(comments_table).
		Join(
			post_comments_table+" on c.id = pc.comment_id",
		).
		Where("id = ?", id)

	if !includeRemoved {
		query = query.Where("removed_at IS NULL")
	}
	sqlQuery, args, _ := query.ToSql()

	dest := struct {
		Id        uuid.UUID    `db:"id"`
		AuthorId  uuid.UUID    `db:"author_id"`
		PostId    uuid.UUID    `db:"post_id"`
		Message   string       `db:"message"`
		CreatedAt time.Time    `db:"created_at"`
		RemovedAt sql.NullTime `db:"removed_at"`
	}{}

	err := this.db.Get(&dest, sqlQuery, args...)
	if err != nil {
		return nil
	}

	builder := domaincomment.NewBuilder().
		Id(dest.Id).
		AuthorId(dest.AuthorId).
		PostId(dest.PostId).
		Message(dest.Message).
		CreatedAt(dest.CreatedAt).
		Permissions(this.permissions)

	if dest.RemovedAt.Valid {
		builder.RemovedAt(dest.RemovedAt.Time)
	}

	return builder.Build()
}

func (this *commentRepository) Save(comment *domaincomment.Comment) *errors.Error {
	qb := postgres.NewSquirrel()
	tx, err := this.db.Beginx()
	if err != nil {
		return errors.From(err)
	}

	for _, event := range comment.Events() {
		switch event.Type() {
		case domaincomment.CREATED:
			{
				sql, args, _ := qb.Insert(comments_table).Columns(
					insert_columns...,
				).Values(
					comment.Id, comment.AuthorId, comment.Message, comment.CreatedAt, comment.RemovedAt,
				).ToSql()
				_, err = tx.Exec(sql, args...)

				if err != nil {
					tx.Rollback()
					return errors.From(err)
				}
			}

			if comment.PostId != uuid.Nil {
				sql, args, _ := qb.Insert("post_comments").Columns("post_id", "comment_id").Values(comment.PostId, comment.Id).ToSql()
				_, err = tx.Exec(sql, args...)

				if err != nil {
					tx.Rollback()
					return errors.From(err)
				}
			}

		case domaincomment.REMOVED:
			sql, args, _ := qb.Update(comments_table).Set("removed_at", comment.RemovedAt).Where("id = ?", comment.Id).ToSql()
			_, err := tx.Exec(sql, args...)
			if err != nil {
				tx.Rollback()
				return errors.From(err)
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.From(err)
	}

	return nil
}
