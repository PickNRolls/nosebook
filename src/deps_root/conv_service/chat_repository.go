package rootconvservice

import (
	domainchat "nosebook/src/domain/chat"
	"nosebook/src/errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type chatRepository struct {
	db *sqlx.DB
}

func newChatRepository(db *sqlx.DB) *chatRepository {
	return &chatRepository{
		db: db,
	}
}

func (this *chatRepository) FindByRecipientId(id uuid.UUID) *domainchat.Chat {
	return nil
}

func (this *chatRepository) Save(chat *domainchat.Chat) *errors.Error {
	return nil
}
