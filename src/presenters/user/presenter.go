package presenteruser

import (
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	presenterdto "nosebook/src/presenters/dto"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Presenter {
	return &Presenter{
		db: db,
	}
}

func (this *Presenter) FindByIds(ids uuid.UUIDs) ([]*presenterdto.User, *errors.Error) {
	qb := postgres.NewSquirrel()
	sql, args, _ := qb.Select(
		"id", "first_name", "last_name", "nick",
	).From(
		"users",
	).Where(
		squirrel.Eq{"id": ids},
	).ToSql()

	output := []*presenterdto.User{}
	err := errors.From(this.db.Select(&output, sql, args...))
	if err != nil {
		return nil, err
	}

	return output, nil
}
