package user

import (
	"context"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"

	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	userRepository UserRepository
	avatarStorage  AvatarStorage
	tracer         trace.Tracer
}

func New(userRepository UserRepository, avatarStorage AvatarStorage, tracer trace.Tracer) *Service {
	return &Service{
		userRepository: userRepository,
		avatarStorage:  avatarStorage,
		tracer:         tracer,
	}
}

func (this *Service) ChangeAvatar(parent context.Context, command ChangeAvatarCommand, auth *auth.Auth) (string, *errors.Error) {
	ctx, span := this.tracer.Start(parent, "user_service.change_avatar")
	defer span.End()

	_, span = this.tracer.Start(ctx, "user_repository.find_by_id")
	user := this.userRepository.FindById(auth.UserId)
	span.End()
	if user == nil {
		return "", errors.New("User Service Error", "Пользователь не найден")
	}

	_, span = this.tracer.Start(ctx, "avatar_storage.upload")
	url, err := this.avatarStorage.Upload(command.Image, auth.UserId)
	span.End()
	if err != nil {
		return "", err
	}

	user.ChangeAvatar(url)

	_, span = this.tracer.Start(ctx, "user_repository.save")
	err = this.userRepository.Save(user)
	span.End()
	if err != nil {
		return "", err
	}

	return url, nil
}
