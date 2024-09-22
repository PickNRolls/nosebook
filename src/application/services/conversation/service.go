package conversation

import (
	"context"
	"nosebook/src/application/services/auth"
	domainchat "nosebook/src/domain/chat"
	"nosebook/src/errors"
	"nosebook/src/lib/clock"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	chatRepository ChatRepository
	userRepository UserRepository
	notifier       Notifier
	tracer         trace.Tracer
}

func New(
	chatRepository ChatRepository,
	notifier Notifier,
	userRepository UserRepository,
	tracer trace.Tracer,
) *Service {
	return &Service{
		chatRepository: chatRepository,
		notifier:       notifier,
		userRepository: userRepository,
		tracer:         tracer,
	}
}

func (this *Service) SendMessage(parent context.Context, command SendMessageCommand, auth *auth.Auth) (bool, *errors.Error) {
  ctx, span := this.tracer.Start(parent, "conversation.send_message") 
  defer span.End()
  
	var err *errors.Error

  _, span = this.tracer.Start(ctx, "user_repository.exists")
	exists := this.userRepository.Exists(command.RecipientId) 
  span.End()
  if !exists {
		return false, newError("Пользователь с id:" + command.RecipientId.String() + " отсуствует")
	}

  _, span = this.tracer.Start(ctx, "chat_repository.find_by_member_ids")
	chat, err := this.chatRepository.FindByMemberIds(command.RecipientId, auth.UserId)
  span.End()
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

  _, span = this.tracer.Start(ctx, "chat_repository.save")
	err = this.chatRepository.Save(chat)
  span.End()
	if err != nil {
		return false, err
	}

  _, span = this.tracer.Start(ctx, "notifier.notify_sender")
	err = this.notifier.NotifyAbout(auth.UserId, chat)
  span.End()
	if err != nil {
		return false, err
	}

  _, span = this.tracer.Start(ctx, "notifier.notifier_recipient")
	err = this.notifier.NotifyAbout(command.RecipientId, chat)
  span.End()
	if err != nil {
		return false, err
	}

	return true, nil
}
