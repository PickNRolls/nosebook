package presenterchat

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type user = presenterdto.User
type message = presenterdto.Message
type chat = presenterdto.Chat
type conversation = presenterdto.Conversation

type MessagePresenter interface {
	FindByIds(ids uuid.UUIDs) (map[uuid.UUID]*message, *errors.Error)
}

type UserPresenter interface {
	FindByIds(ids uuid.UUIDs) (map[uuid.UUID]*user, *errors.Error)
}
