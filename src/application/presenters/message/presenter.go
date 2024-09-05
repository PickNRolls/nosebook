package presentermessage

import (
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/nullable"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type dest struct {
	Id        uuid.UUID     `db:"id"`
	AuthorId  uuid.UUID     `db:"author_id"`
	Text      string        `db:"text"`
	ReplyTo   nullable.UUID `db:"reply_to"`
	CreatedAt time.Time     `db:"created_at"`
}

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

func (this *Presenter) FindByIds(ids uuid.UUIDs) (map[uuid.UUID]*message, *errors.Error) {
	qb := querybuilder.New()
	sql, args, _ := qb.
		Select("id", "author_id", "text", "reply_to", "created_at").
		From("messages").
		Where(squirrel.Eq{"id": ids}).
		ToSql()

	dests := []*dest{}
	err := this.db.Select(&dests, sql, args...)
	if err != nil {
		return nil, errors.From(err)
	}

	userMap, err := func() (map[uuid.UUID]*user, *errors.Error) {
		ids := []uuid.UUID{}
		idMap := make(map[uuid.UUID]struct{})

		for _, dest := range dests {
			if _, has := idMap[dest.AuthorId]; !has {
				idMap[dest.AuthorId] = struct{}{}
				ids = append(ids, dest.AuthorId)
			}
		}

		return this.userPresenter.FindByIds(ids)
	}()

	out := make(map[uuid.UUID]*message, len(dests))
	for _, dest := range dests {
		out[dest.Id] = &message{
			Id:        dest.Id,
			Author:    userMap[dest.AuthorId],
			Text:      dest.Text,
			CreatedAt: dest.CreatedAt,
		}

		if dest.ReplyTo.Valid {
			out[dest.Id].ReplyTo = &dest.ReplyTo
		}
	}
	return out, nil
}
