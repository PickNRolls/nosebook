package presenteruser

import (
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	"nosebook/src/lib/clock"
	"time"

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

func (this *Presenter) FindByIds(ids uuid.UUIDs) ([]*User, *errors.Error) {
	qb := postgres.NewSquirrel()
	sql, args, _ := qb.Select(
		"id", "first_name", "last_name", "nick", "last_activity_at",
	).From(
		"users",
	).Where(
		squirrel.Eq{"id": ids},
	).ToSql()

	dest := []*dest{}
	err := errors.From(this.db.Select(&dest, sql, args...))
	if err != nil {
		return nil, err
	}

	output := func() []*User {
		out := make([]*User, len(dest))

		now := clock.Now()

		for i, userDest := range dest {
			out[i] = &User{
				Id:           userDest.Id,
				FirstName:    userDest.FirstName,
				LastName:     userDest.LastName,
				Nick:         userDest.Nick,
				LastOnlineAt: userDest.LastActivityAt,
				Online:       userDest.LastActivityAt.After(now.Add(-5 * time.Minute)),
			}
		}

		return out
	}()

	return output, nil
}
