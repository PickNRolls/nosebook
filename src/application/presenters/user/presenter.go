package presenteruser

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	domainuser "nosebook/src/domain/user"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/clock"
	cursorquery "nosebook/src/lib/cursor_query"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Presenter struct {
	db     *sqlx.DB
	tracer trace.Tracer
}

func New(db *sqlx.DB) *Presenter {
	return &Presenter{
		db:     db,
		tracer: noop.Tracer{},
	}
}

func (this *Presenter) WithTracer(tracer trace.Tracer) *Presenter {
	this.tracer = tracer

	return this
}

func (this *Presenter) mapDests(dests []*dest) map[uuid.UUID]*User {
	out := make(map[uuid.UUID]*User, len(dests))

	now := clock.Now()
	for _, dest := range dests {
		out[dest.Id] = &User{
			Id:           dest.Id,
			FirstName:    dest.FirstName,
			LastName:     dest.LastName,
			Nick:         dest.Nick,
			LastOnlineAt: dest.LastActivityAt,
			Online:       dest.LastActivityAt.After(now.Add(-domainuser.ONLINE_DURATION)),
		}

		if dest.AvatarUrl != "" {
			out[dest.Id].Avatar = &presenterdto.UserAvatar{
				Url:       dest.AvatarUrl,
				UpdatedAt: *dest.AvatarUpdatedAt,
			}
		}
	}

	return out
}

func (this *Presenter) FindByIds(ctx context.Context, ids uuid.UUIDs) (map[uuid.UUID]*User, *errors.Error) {
	nextCtx, span := this.tracer.Start(ctx, "user_presenter.find_by_filter")
	defer span.End()

	qb := querybuilder.New()
	sql, args, _ := qb.Select(
		"id", "first_name", "last_name", "nick", "avatar_url", "avatar_updated_at", "last_activity_at",
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

	return this.mapDests(dests), nil
}

func outErr(err error) *FindOutUser {
	return outMsgErr(err.Error())
}

func outMsgErr(message string) *FindOutUser {
	return &FindOutUser{
		Err: newError(message),
	}
}

func outZero() *FindOutUser {
	return &FindOutUser{
		Data: make([]*User, 0),
	}
}

func (this *Presenter) FindByText(parent context.Context, input FindByTextInput, _ *auth.Auth) *FindOutUser {
	ctx, span := this.tracer.Start(parent, "user_presenter.find_by_text")
	defer span.End()

	if input.Text == "" {
		return outZero()
	}

	qb := querybuilder.New()
	query := qb.Select(
		"id", "first_name", "last_name", "nick", "avatar_url", "avatar_updated_at", "last_activity_at", "created_at",
	).From(
		"users",
	).Where(
		"(first_name || ' ' || last_name || ' ' || nick) ilike '%' || ? || '%'",
		input.Text,
	)

	dests := []*dest{}
	_, span = this.tracer.Start(ctx, "sql_query")
	cursorOut, err := cursorquery.Do(this.db, &cursorquery.Input[*dest]{
		Query: query,
		Order: &order{},
		Next:  input.Next,
		Limit: 20,
	}, &dests)
	span.End()
	if err != nil {
		return outErr(err)
	}

	users := func() []*User {
		m := this.mapDests(dests)
		out := make([]*User, len(dests))

		for i, dest := range dests {
			out[i] = m[dest.Id]
		}

		return out
	}()

	return &FindOutUser{
		Data:       users,
		Next:       cursorOut.Next,
		TotalCount: cursorOut.TotalCount,
	}
}
