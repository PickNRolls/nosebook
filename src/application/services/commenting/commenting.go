package commenting

import (
	"context"
	"nosebook/src/application/services/auth"
	"nosebook/src/domain/comment"
	"nosebook/src/errors"
	"nosebook/src/lib/clock"

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

func (this *Service) PublishOnPost(parent context.Context, c PublishPostCommentCommand, a *auth.Auth) (uuid.UUID, *errors.Error) {
	if post := this.postRepository.FindById(c.Id); post == nil {
		return uuid.Nil, NewPostNotFoundError()
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
		return uuid.Nil, err
	}

	return comment.Id, nil
}

func (this *Service) Remove(parent context.Context, c RemoveCommentCommand, a *auth.Auth) (uuid.UUID, *errors.Error) {
	comment := this.repository.FindById(c.Id, true)
	if comment == nil {
    return uuid.Nil, NewError("Такого комментария не существует")
	}

	err := comment.RemoveBy(a.UserId)
	if err != nil {
    return uuid.Nil, err
	}

	err = this.repository.Save(comment)
	if err != nil {
    return uuid.Nil, err
	}

  return comment.Id, nil
}
