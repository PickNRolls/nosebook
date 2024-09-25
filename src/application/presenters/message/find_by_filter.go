package presentermessage

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	cursorquery "nosebook/src/lib/cursor_query"
	"time"

	"github.com/google/uuid"
)

type order struct{}

func (this *order) Column() string {
	return "created_at"
}
func (this *order) Timestamp(dest *dest) time.Time {
	return dest.CreatedAt
}
func (this *order) Id(dest *dest) uuid.UUID {
	return dest.Id
}
func (this *order) Asc() bool {
	return false
}

type FindByFilterInput struct {
	ChatId string
	Next   string
	Prev   string
	Limit  uint64
}

type FindByFilterOut presenterdto.FindOut[*message]

func errMsgOut(message string) *FindByFilterOut {
	return &FindByFilterOut{
		Err: errors.New("Message Presenter Error", message),
	}
}

func errOut(err error) *FindByFilterOut {
	return errMsgOut(err.Error())
}

func (this *Presenter) FindByFilter(parent context.Context, input FindByFilterInput, auth *auth.Auth) *FindByFilterOut {
  ctx, span := this.tracer.Start(parent, "message_presenter.find_by_filter")
  defer span.End()
  
	if input.ChatId == "" {
		return errMsgOut("ChatId - обязательный параметр")
	}

	chatId, err := errors.Using(uuid.Parse(input.ChatId))
	if err != nil {
		return errOut(err)
	}

	qb := querybuilder.New()
	query := qb.
		Select("id", "author_id", "text", "chat_id", "reply_to", "created_at").
		From("messages").
		Where("chat_id = ?", chatId)

  _, span = this.tracer.Start(ctx, "message_presenter.sql_query")
	dests := []*dest{}
	cursorQueryOut, err := cursorquery.Do(this.db, &cursorquery.Input[*dest]{
		Query: query,
		Order: &order{},
		Next:  input.Next,
		Prev:  input.Prev,
		Limit: input.Limit,
	}, &dests)
  span.End()
	if err != nil {
		return errOut(err)
	}

	messageMap, err := mapDests(ctx, this.userPresenter, dests)
	if err != nil {
		return errOut(err)
	}

	out := &FindByFilterOut{
		Data:       make([]*presenterdto.Message, len(dests)),
		TotalCount: cursorQueryOut.TotalCount,
		Next:       cursorQueryOut.Next,
		Prev:       cursorQueryOut.Prev,
	}
	for i, dest := range dests {
		out.Data[i] = messageMap[dest.Id]
	}
	return out
}
