package presenterchat

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	cursorquery "nosebook/src/lib/cursor_query"

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

func (this *Presenter) FindByFilter(input *FindByFilterInput, auth *auth.Auth) *FindByFilterOutput {
	qb := querybuilder.New()
	query := qb.Select("id", "created_at", "interlocutor_id", "last_message_id").
		FromSelect(qb.Select(
			"c.id",
			"c.created_at",
			"cm2.user_id as interlocutor_id",
			"m.id as last_message_id",
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
	cursorQueryOut, err := cursorquery.Do(this.db, &cursorquery.Input{
		Query: query,
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
