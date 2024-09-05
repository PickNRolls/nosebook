package conversation

import (
	domainchat "nosebook/src/domain/chat"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type ChatRepository interface {
	FindByRecipientId(id uuid.UUID) *domainchat.Chat
	Save(chat *domainchat.Chat) *errors.Error
}

type UserRepository interface {
	Exists(id uuid.UUID) bool
}

type Notifier interface {
	Notify(chat *domainchat.Chat)
}

type NotifierRepository interface {
	FindByRecipientId(id uuid.UUID) Notifier
}
