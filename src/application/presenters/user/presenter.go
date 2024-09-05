package presenteruser

import (
	domainuser "nosebook/src/domain/user"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/clock"

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

func (this *Presenter) FindByIds(ids uuid.UUIDs) (map[uuid.UUID]*User, *errors.Error) {
	qb := querybuilder.New()
	sql, args, _ := qb.Select(
		"id", "first_name", "last_name", "nick", "last_activity_at",
	).From(
		"users",
	).Where(
		squirrel.Eq{"id": ids},
	).ToSql()

	dests := []*dest{}
	err := errors.From(this.db.Select(&dests, sql, args...))
	if err != nil {
		return nil, err
	}

	m := make(map[uuid.UUID]*User, len(dests))
	now := clock.Now()
	for _, dest := range dests {
		m[dest.Id] = &User{
			Id:           dest.Id,
			FirstName:    dest.FirstName,
			LastName:     dest.LastName,
			Nick:         dest.Nick,
			LastOnlineAt: dest.LastActivityAt,
			Online:       dest.LastActivityAt.After(now.Add(-domainuser.ONLINE_DURATION)),
		}
	}
	return m, nil
}
