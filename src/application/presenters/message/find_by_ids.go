package presentermessage

import (
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (this *Presenter) FindByIds(ids uuid.UUIDs) (map[uuid.UUID]*message, *errors.Error) {
	qb := querybuilder.New()
	sql, args, _ := qb.
		Select("id", "author_id", "text", "chat_id", "reply_to", "created_at").
		From("messages").
		Where(squirrel.Eq{"id": ids}).
		ToSql()

	dests := []*dest{}
	err := this.db.Select(&dests, sql, args...)
	if err != nil {
		return nil, errors.From(err)
	}

	return this.mapDests(dests)
}
