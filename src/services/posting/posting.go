package posting

import (
	"nosebook/src/domain/post"
	"nosebook/src/services/auth"
	commandresult "nosebook/src/services/command_result"
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

func (this *Service) Publish(c *PublishPostCommand, a *auth.Auth) *commandresult.Result {
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
		return commandresult.Fail(err)
	}

	return commandresult.Ok().WithId(post.Id)
}

func (this *Service) Remove(c *RemovePostCommand, a *auth.Auth) *commandresult.Result {
	post := this.repository.FindById(c.Id)
	if post == nil {
		return commandresult.Fail(NewNotFoundError())
	}

	err := post.RemoveBy(a.UserId)
	if err != nil {
		return commandresult.Fail(err)
	}

	err = this.repository.Save(post)
	if err != nil {
		return commandresult.Fail(err)
	}

	return commandresult.Ok().WithId(post.Id)
}

func (this *Service) Edit(c *EditPostCommand, a *auth.Auth) *commandresult.Result {
	post := this.repository.FindById(c.Id)
	if post == nil {
		return commandresult.Fail(NewNotFoundError())
	}

	err := post.EditBy(a.UserId, c.Message)
	if err != nil {
		return commandresult.Fail(err)
	}

	err = this.repository.Save(post)
	if err != nil {
		return commandresult.Fail(err)
	}

	return commandresult.Ok().WithId(post.Id)
}
