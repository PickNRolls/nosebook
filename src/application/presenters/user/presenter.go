package presenteruser

import (
	"context"
	domainuser "nosebook/src/domain/user"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/clock"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type startCallback = func(name string, ctx context.Context) func ()

type Presenter struct {
	db      *sqlx.DB
  tracer trace.Tracer
}

func New(db *sqlx.DB) *Presenter {
	return &Presenter{
		db: db,
    tracer: noop.Tracer{},
	}
}

func (this *Presenter) WithTracer(tracer trace.Tracer) *Presenter {
  this.tracer = tracer

  return this
}

func (this *Presenter) FindByIds(ctx context.Context, ids uuid.UUIDs) (map[uuid.UUID]*User, *errors.Error) {
  nextCtx, span := this.tracer.Start(ctx, "user_presenter.find_by_filter")
  defer span.End()

	qb := querybuilder.New()
	sql, args, _ := qb.Select(
		"id", "first_name", "last_name", "nick", "last_activity_at",
	).From(
		"users",
	).Where(
		squirrel.Eq{"id": ids},
	).ToSql()

	dests := []*dest{}
  nextCtx, span = this.tracer.Start(nextCtx, "user_presenter.sql_query")
	err := errors.From(this.db.Select(&dests, sql, args...))
  span.End()
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
