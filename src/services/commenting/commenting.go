package commenting

import (
	"nosebook/src/domain/comment"
	"nosebook/src/errors"
	"nosebook/src/services/auth"
	"time"

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

func (this *Service) PublishOnPost(c *PublishPostCommentCommand, a *auth.Auth) (*domaincomment.Comment, *errors.Error) {
	if post := this.postRepository.FindById(c.Id); post == nil {
		return nil, NewPostNotFoundError()
	}

	comment := domaincomment.NewBuilder().
		Id(uuid.New()).
		AuthorId(a.UserId).
		Message(c.Message).
		PostId(c.Id).
		CreatedAt(time.Now()).
		RaiseCreatedEvent().
		Build()

	err := this.repository.Save(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (this *Service) Remove(c *RemoveCommentCommand, a *auth.Auth) (*domaincomment.Comment, *errors.Error) {
	comment := this.repository.FindById(c.Id, true)
	if comment == nil {
		return nil, NewError("Такого комментария не существует")
	}

	err := comment.RemoveBy(a.UserId)
	if err != nil {
		return nil, err
	}

	err = this.repository.Save(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
