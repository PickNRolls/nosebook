package conversation

import (
	domainchat "nosebook/src/domain/chat"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type ChatRepository interface {
	FindByMemberIds(leftId uuid.UUID, rightId uuid.UUID) (*domainchat.Chat, *errors.Error)
	Save(chat *domainchat.Chat) *errors.Error
}

type UserRepository interface {
	Exists(id uuid.UUID) bool
}

type Notifier interface {
	Notify(chat *domainchat.Chat) *errors.Error
}

type NotifierRepository interface {
	FindByUserId(id uuid.UUID) Notifier
}
