package presenterchat

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	cursorquery "nosebook/src/lib/cursor_query"
	"time"

	"github.com/google/uuid"
)

type FindByFilterInput struct {
	Next  string
	Limit uint64
}

type FindByFilterOutput presenterdto.FindOut[chat]

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

func (this *Presenter) FindByFilter(input *FindByFilterInput, auth *auth.Auth) *FindByFilterOutput {
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

	dests := []*conv_dest{}
	cursorQueryOut, err := cursorquery.Do(this.db, &cursorquery.Input[*conv_dest]{
		Query: query,
		Order: &order{},
		Next:  input.Next,
		Limit: input.Limit,
	}, &dests)
	if err != nil {
		return errOut(err)
	}

	userMap, err := func() (map[uuid.UUID]*user, *errors.Error) {
		ids := make([]uuid.UUID, len(dests))
		for i, dest := range dests {
			ids[i] = dest.InterlocutorId
		}

		return this.userPresenter.FindByIds(ids)
	}()
	if err != nil {
		return errOut(err)
	}

	messageMap, err := func() (map[uuid.UUID]*message, *errors.Error) {
		ids := make([]uuid.UUID, len(dests))
		for i, dest := range dests {
			ids[i] = dest.LastMessageId
		}

		return this.messagePresenter.FindByIds(ids)
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
