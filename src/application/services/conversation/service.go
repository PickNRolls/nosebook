package conversation

import (
	"nosebook/src/application/services/auth"
	domainchat "nosebook/src/domain/chat"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type Service struct {
	chatRepository     ChatRepository
	userRepository     UserRepository
	notifierRepository NotifierRepository
}

func New(
	chatRepository ChatRepository,
	notifierRepository NotifierRepository,
	userRepository UserRepository,
) *Service {
	return &Service{
		chatRepository:     chatRepository,
		notifierRepository: notifierRepository,
		userRepository:     userRepository,
	}
}

func (this *Service) SendMessage(command *SendMessageCommand, auth *auth.Auth) (bool, *errors.Error) {
	var err *errors.Error

	if exists := this.userRepository.Exists(command.RecipientId); !exists {
		return false, newError("Пользователь с id:" + command.RecipientId.String() + " отсуствует")
	}

	chat := this.chatRepository.FindByRecipientId(command.RecipientId)
	if chat == nil {
		chat, err = domainchat.New(
			uuid.New(),
			[]uuid.UUID{auth.UserId, command.RecipientId},
			"",
			true,
			nil,
			true,
		)
	}

	if err != nil {
		return false, err
	}

	err = chat.SendMessageBy(command.Text, command.ReplyTo, auth.UserId)
	if err != nil {
		return false, err
	}

	err = this.chatRepository.Save(chat)
	if err != nil {
		return false, err
	}

	notifier := this.notifierRepository.FindByRecipientId(command.RecipientId)
	if notifier == nil {
		return true, nil
	}

	notifier.Notify(chat)

	return true, nil
}
