package presenterchat

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type user = presenterdto.User
type message = presenterdto.Message
type chat = presenterdto.Chat
type conversation = presenterdto.Conversation

type MessagePresenter interface {
	FindByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]*message, *errors.Error)
}

type UserPresenter interface {
	FindByIds(ctx context.Context, ids uuid.UUIDs) (map[uuid.UUID]*user, *errors.Error)
}
