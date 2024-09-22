package posting

import (
	"context"
	"nosebook/src/application/services/auth"
	"nosebook/src/domain/post"
	"nosebook/src/errors"
	"nosebook/src/lib/clock"

	"github.com/google/uuid"
)

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (this *Service) Publish(parent context.Context, c PublishPostCommand, a *auth.Auth) (uuid.UUID, *errors.Error) {
	post := domainpost.NewBuilder().
		Id(uuid.New()).
		AuthorId(a.UserId).
		OwnerId(c.OwnerId).
		Message(c.Message).
		CreatedAt(clock.Now()).
		RaiseCreatedEvent().
		Build()

	err := this.repository.Save(post)
	if err != nil {
		return uuid.Nil, err
	}

	return post.Id, nil
}

func (this *Service) Remove(parent context.Context, c RemovePostCommand, a *auth.Auth) (uuid.UUID, *errors.Error) {
	post := this.repository.FindById(c.Id)
	if post == nil {
		return uuid.Nil, NewNotFoundError()
	}

	err := post.RemoveBy(a.UserId)
	if err != nil {
		return uuid.Nil, err
	}

	err = this.repository.Save(post)
	if err != nil {
		return uuid.Nil, err
	}

	return post.Id, nil
}

func (this *Service) Edit(c EditPostCommand, a *auth.Auth) (uuid.UUID, *errors.Error) {
	post := this.repository.FindById(c.Id)
	if post == nil {
		return uuid.Nil, NewNotFoundError()
	}

	err := post.EditBy(a.UserId, c.Message)
	if err != nil {
		return uuid.Nil, err
	}

	err = this.repository.Save(post)
	if err != nil {
		return uuid.Nil, err
	}

	return post.Id, nil
}
