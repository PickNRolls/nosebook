package presenterfriendship

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	domainuser "nosebook/src/domain/user"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/clock"
	cursorquery "nosebook/src/lib/cursor_query"
	"nosebook/src/lib/nullable"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type FindByFilterInput struct {
	UserId string

	Accepted nullable.Bool
	Viewed   nullable.Bool

	OnlyIncoming  bool
	OnlyOutcoming bool
	OnlyOnline    bool

	Next  string
	Prev  string
	Limit uint64
}

type FindByFilterOutput = presenterdto.FindOut[*Request]

func errMsgOut(message string) *FindByFilterOutput {
	return &FindByFilterOutput{
		Err: newError(message),
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

	incomingQuery := squirrel.StatementBuilder.
		Select("requester_id as id, created_at, accepted, 'incoming' as type").
		From("friendship_requests").
		Where("responder_id = ?", userId)

	outcomingQuery := squirrel.StatementBuilder.
		Select("responder_id as id, created_at, accepted, 'outcoming' as type").
		From("friendship_requests").
		Where("requester_id = ?", userId)

	if input.Accepted.Valid {
		incomingQuery = incomingQuery.Where("accepted = ?", input.Accepted.Value)
		outcomingQuery = outcomingQuery.Where("accepted = ?", input.Accepted.Value)
	}

	if input.Viewed.Valid {
		incomingQuery = incomingQuery.Where("viewed = ?", input.Viewed.Value)
		outcomingQuery = outcomingQuery.Where("viewed = ?", input.Viewed.Value)
	}

	union := querybuilder.Union(
		incomingQuery,
		outcomingQuery,
	).PlaceholderFormat(squirrel.Dollar)

	innerQuery := union
	if input.OnlyIncoming {
		innerQuery = incomingQuery
	}

	if input.OnlyOutcoming {
		innerQuery = outcomingQuery
	}

	query := querybuilder.New().
		Select("f.id", "f.created_at", "f.accepted", "f.type").
		FromSelect(innerQuery, "f")

	if input.OnlyOnline {
		query = query.
			Join("users as u on u.id = f.id").
			Where("u.last_activity_at > ?", clock.Now().Add(-domainuser.ONLINE_DURATION))
	}

	dests := []*find_by_filter_dest{}
	cursorQueryOut, error := cursorquery.Do(this.db, &cursorquery.Input[*find_by_filter_dest]{
		Query: query,
		Order: &find_by_filter_order{},
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

		return this.userPresenter.FindByIds(context.TODO(), ids)
	}()

	output := &FindByFilterOutput{}
	output.Data = make([]*Request, len(dests))
	for i, dest := range dests {
		output.Data[i] = &Request{
			Type:     dest.Type,
			Accepted: dest.Accepted,
			User:     userMap[dest.Id],
		}
	}
	output.Next = cursorQueryOut.Next
	output.Prev = cursorQueryOut.Prev
	output.TotalCount = cursorQueryOut.TotalCount

	return output
}
