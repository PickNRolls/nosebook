package posting

import (
	"nosebook/src/domain/post"
	"nosebook/src/errors"
	"nosebook/src/services/auth"
	"time"

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

func (this *Service) Publish(c *PublishPostCommand, a *auth.Auth) (*domainpost.Post, *errors.Error) {
	post := domainpost.NewBuilder().
		Id(uuid.New()).
		AuthorId(a.UserId).
		OwnerId(c.OwnerId).
		Message(c.Message).
		CreatedAt(time.Now()).
		RaiseCreatedEvent().
		Build()

	err := this.repository.Save(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (this *Service) Remove(c *RemovePostCommand, a *auth.Auth) (*domainpost.Post, *errors.Error) {
	post := this.repository.FindById(c.Id)
	if post == nil {
		return nil, NewNotFoundError()
	}

	err := post.RemoveBy(a.UserId)
	if err != nil {
		return nil, err
	}

	err = this.repository.Save(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (this *Service) Edit(c *EditPostCommand, a *auth.Auth) (*domainpost.Post, *errors.Error) {
	post := this.repository.FindById(c.Id)
	if post == nil {
		return nil, NewNotFoundError()
	}

	err := post.EditBy(a.UserId, c.Message)
	if err != nil {
		return nil, err
	}

	err = this.repository.Save(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}
