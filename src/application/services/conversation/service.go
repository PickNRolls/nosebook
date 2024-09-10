package conversation

import (
	"nosebook/src/application/services/auth"
	domainchat "nosebook/src/domain/chat"
	"nosebook/src/errors"
	"nosebook/src/lib/clock"

	"github.com/google/uuid"
)

type Service struct {
	chatRepository ChatRepository
	userRepository UserRepository
	notifier       Notifier
}

func New(
	chatRepository ChatRepository,
	notifier Notifier,
	userRepository UserRepository,
) *Service {
	return &Service{
		chatRepository: chatRepository,
		notifier:       notifier,
		userRepository: userRepository,
	}
}

func (this *Service) SendMessage(command *SendMessageCommand, auth *auth.Auth) (bool, *errors.Error) {
	var err *errors.Error

	if exists := this.userRepository.Exists(command.RecipientId); !exists {
		return false, newError("Пользователь с id:" + command.RecipientId.String() + " отсуствует")
	}

	chat, err := this.chatRepository.FindByMemberIds(command.RecipientId, auth.UserId)
	if err != nil {
		return false, err
	}

	if chat == nil {
		chat, err = domainchat.New(
			uuid.New(),
			[]uuid.UUID{auth.UserId, command.RecipientId},
			"",
			true,
			clock.Now(),
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

	err = this.notifier.NotifyAbout(auth.UserId, chat)
	if err != nil {
		return false, err
	}
	err = this.notifier.NotifyAbout(command.RecipientId, chat)
	if err != nil {
		return false, err
	}

	return true, nil
}
