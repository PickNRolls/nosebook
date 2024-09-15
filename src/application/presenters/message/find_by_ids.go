package presentermessage

import (
	"context"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (this *Presenter) FindByIds(ctx context.Context, ids uuid.UUIDs) (map[uuid.UUID]*message, *errors.Error) {
  nextCtx, span := this.tracer.Start(ctx, "message_presenter.find_by_ids")
  defer span.End()
  
	qb := querybuilder.New()
	sql, args, _ := qb.
		Select("id", "author_id", "text", "chat_id", "reply_to", "created_at").
		From("messages").
		Where(squirrel.Eq{"id": ids}).
		ToSql()

	dests := []*dest{}
  _, span = this.tracer.Start(nextCtx, "message_presenter.sql_query")
	err := this.db.Select(&dests, sql, args...)
  span.End()
	if err != nil {
		return nil, errors.From(err)
	}

	return this.mapDests(nextCtx, dests)
}
