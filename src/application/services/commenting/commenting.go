package commenting

import (
	"nosebook/src/application/services/auth"
	"nosebook/src/domain/comment"
	"nosebook/src/lib/clock"
	commandresult "nosebook/src/lib/command_result"

	"github.com/google/uuid"
)

type Service struct {
	repository     Repository
	postRepository PostRepository
}

func New(repository Repository, postRepository PostRepository) *Service {
	return &Service{
		repository:     repository,
		postRepository: postRepository,
	}
}

func (this *Service) PublishOnPost(c *PublishPostCommentCommand, a *auth.Auth) *commandresult.Result {
	if post := this.postRepository.FindById(c.Id); post == nil {
		return commandresult.Fail(NewPostNotFoundError())
	}

	comment := domaincomment.NewBuilder().
		Id(uuid.New()).
		AuthorId(a.UserId).
		Message(c.Message).
		PostId(c.Id).
		CreatedAt(clock.Now()).
		RaiseCreatedEvent().
		Build()

	err := this.repository.Save(comment)
	if err != nil {
		return commandresult.Fail(err)
	}

	return commandresult.Ok().WithId(comment.Id)
}

func (this *Service) Remove(c *RemoveCommentCommand, a *auth.Auth) *commandresult.Result {
	comment := this.repository.FindById(c.Id, true)
	if comment == nil {
		return commandresult.Fail(NewError("Такого комментария не существует"))
	}

	err := comment.RemoveBy(a.UserId)
	if err != nil {
		return commandresult.Fail(err)
	}

	err = this.repository.Save(comment)
	if err != nil {
		return commandresult.Fail(err)
	}

	return commandresult.Ok().WithId(comment.Id)
}
