package presentermessage

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/errors"
	"nosebook/src/lib/nullable"
	"time"

	"github.com/google/uuid"
)

type dest struct {
	Id        uuid.UUID     `db:"id"`
	AuthorId  uuid.UUID     `db:"author_id"`
	Text      string        `db:"text"`
	ReplyTo   nullable.UUID `db:"reply_to"`
	ChatId    uuid.UUID     `db:"chat_id"`
	CreatedAt time.Time     `db:"created_at"`
}

type message = presenterdto.Message
type user = presenterdto.User

type UserPresenter interface {
	FindByIds(ids uuid.UUIDs) (map[uuid.UUID]*user, *errors.Error)
}
