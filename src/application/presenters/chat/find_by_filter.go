package presenterchat

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	cursorquery "nosebook/src/lib/cursor_query"
	"nosebook/src/lib/nullable"
	"time"

	"github.com/google/uuid"
)

type FindByFilterInput struct {
	Id             string
	InterlocutorId string

	Next  string
	Limit uint64
}

func (this *FindByFilterInput) BuildFromMap(m map[string]any) FindByFilterInput {
	out := FindByFilterInput{}

	if id, ok := m["id"].(string); ok {
		out.Id = id
	}

	if interlocutorId, ok := m["interlocutorId"].(string); ok {
		out.InterlocutorId = interlocutorId
	}

	if next, ok := m["next"].(string); ok {
		out.Next = next
	}

	if limit, ok := m["limit"].(nullable.Uint64); ok && limit.Valid {
		out.Limit = limit.Value
	}

	return out
}

type FindByFilterOutput = presenterdto.FindOut[chat]

func errOut(err error) *FindByFilterOutput {
	return &FindByFilterOutput{
		Err: errors.From(err),
	}
}

type order struct{}

func (this *order) Column() string {
	return "updated_at"
}
func (this *order) Timestamp(dest *conv_dest) time.Time {
	return dest.UpdatedAt
}
func (this *order) Id(dest *conv_dest) uuid.UUID {
	return dest.Id
}
func (this *order) Asc() bool {
	return false
}

func (this *Presenter) FindByFilter(ctx context.Context, input FindByFilterInput, auth *auth.Auth) *FindByFilterOutput {
	nextCtx, span := this.tracer.Start(ctx, "chat_presenter.find_by_filter")
	defer span.End()

	var id uuid.UUID
	var err *errors.Error
	if input.Id != "" {
		id, err = errors.Using(uuid.Parse(input.Id))
		if err != nil {
			return errOut(err)
		}
	}

	var interlocutorId uuid.UUID
	if input.InterlocutorId != "" {
		interlocutorId, err = errors.Using(uuid.Parse(input.InterlocutorId))
		if err != nil {
			return errOut(err)
		}
	}

	qb := querybuilder.New()
	query := qb.Select(
		"id",
		"interlocutor_id",
		"last_message_id",
		"created_at",
		"updated_at",
	).
		FromSelect(qb.Select(
			"c.id",
			"c.created_at",
			"cm2.user_id as interlocutor_id",
			"m.id as last_message_id",
			"m.created_at as updated_at",
			"row_number() over (partition by c.id order by m.created_at desc) as row_number",
		).
			From("chats as c").
			Join("chat_members as cm on c.id = cm.chat_id").
			Join("chat_members as cm2 on cm.chat_id = cm2.chat_id AND cm.user_id != cm2.user_id").
			Join("messages as m on m.chat_id = c.id").
			Where("cm.user_id = ?", auth.UserId).
			OrderBy("m.created_at desc"),
			"c",
		).Where("c.row_number <= 1")

	if interlocutorId != uuid.Nil {
		query = query.Where("interlocutor_id = ?", interlocutorId)
	}

	if id != uuid.Nil {
		query = query.Where("id = ?", id)
	}

	dests := []*conv_dest{}
  
	_, span = this.tracer.Start(nextCtx, "chat_presenter.sql_query")
	cursorQueryOut, err := cursorquery.Do(this.db, &cursorquery.Input[*conv_dest]{
		Query: query,
		Order: &order{},
		Next:  input.Next,
		Limit: input.Limit,
	}, &dests)
  span.End()
	if err != nil {
		return errOut(err)
	}

	userMap, err := func() (map[uuid.UUID]*user, *errors.Error) {
		ids := make([]uuid.UUID, len(dests))
		for i, dest := range dests {
			ids[i] = dest.InterlocutorId
		}

		return this.userPresenter.FindByIds(nextCtx, ids)
	}()
	if err != nil {
		return errOut(err)
	}

	messageMap, err := func() (map[uuid.UUID]*message, *errors.Error) {
		ids := make([]uuid.UUID, len(dests))
		for i, dest := range dests {
			ids[i] = dest.LastMessageId
		}

		return this.messagePresenter.FindByIds(nextCtx, ids)
	}()
	if err != nil {
		return errOut(err)
	}

	data := make([]chat, len(dests))
	for i, dest := range dests {
		data[i] = &conversation{
			Id:           dest.Id,
			Interlocutor: userMap[dest.InterlocutorId],
			LastMessage:  messageMap[dest.LastMessageId],
			CreatedAt:    dest.CreatedAt,
		}
	}

	return &FindByFilterOutput{
		Data:       data,
		Next:       cursorQueryOut.Next,
		TotalCount: cursorQueryOut.TotalCount,
	}
}
