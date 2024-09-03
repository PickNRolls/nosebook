package presenterfriendship

import (
	"nosebook/src/application/services/auth"
	domainuser "nosebook/src/domain/user"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/clock"
	cursorquery "nosebook/src/lib/cursor_query"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	db            *sqlx.DB
	userPresenter UserPresenter
}

func New(db *sqlx.DB, userPresenter UserPresenter) *Presenter {
	return &Presenter{
		db:            db,
		userPresenter: userPresenter,
	}
}

func errMsgOut(message string) *FindByFilterOutput {
	return &FindByFilterOutput{
		Err: errors.New("Friendship Presenter Error", message),
	}
}

func errOut(err error) *FindByFilterOutput {
	return errMsgOut(err.Error())
}

func (this *Presenter) FindByFilter(input *FindByFilterInput, auth *auth.Auth) *FindByFilterOutput {
	if input.UserId == "" {
		return errMsgOut("UserId - обязательный параметр")
	}

	userId, err := uuid.Parse(input.UserId)
	if err != nil {
		return errOut(err)
	}

	union := querybuilder.Union(
		squirrel.StatementBuilder.
			Select("requester_id as id, created_at").
			From("friendship_requests").
			Where("accepted = true").
			Where("responder_id = ?", userId),

		squirrel.StatementBuilder.
			Select("responder_id as id, created_at").
			From("friendship_requests").
			Where("accepted = true").
			Where("requester_id = ?", userId),
	).PlaceholderFormat(squirrel.Dollar)

	query := querybuilder.New().
		Select("f.id", "f.created_at").
		FromSelect(union, "f")

	if input.OnlyOnline {
		query = query.
			Join("users as u on u.id = f.id").
			Where("u.last_activity_at > ?", clock.Now().Add(-domainuser.ONLINE_DURATION))
	}

	dests := []*dest{}
	cursorQueryOut, error := cursorquery.Do(this.db, &cursorquery.Input{
		Query: query,
		Next:  input.Next,
		Limit: input.Limit,
	}, &dests)

	if error != nil {
		return errOut(error)
	}

	userMap, error := func() (map[uuid.UUID]*user, *errors.Error) {
		ids := make([]uuid.UUID, len(dests))
		for i, dest := range dests {
			ids[i] = dest.Id
		}

		users, err := this.userPresenter.FindByIds(ids)
		if err != nil {
			return nil, err
		}

		m := make(map[uuid.UUID]*user, len(users))
		for _, user := range users {
			m[user.Id] = user
		}
		return m, nil
	}()

	output := &FindByFilterOutput{}
	output.Data = make([]*user, len(dests))
	for i, dest := range dests {
		output.Data[i] = userMap[dest.Id]
	}
	output.Next = cursorQueryOut.Next
	output.Prev = cursorQueryOut.Prev
	output.TotalCount = cursorQueryOut.TotalCount

	return output
}
