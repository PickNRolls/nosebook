package services

import (
	"nosebook/src/domain/comments"
	"nosebook/src/errors"
	"nosebook/src/services/auth"
	"nosebook/src/services/commenting"
	"nosebook/src/services/commenting/commands"
	"nosebook/src/services/commenting/interfaces"
	"time"

	"github.com/google/uuid"
)

type CommentService struct {
	repository     interfaces.CommentRepository
	postRepository interfaces.PostRepository
}

func NewCommentService(repository interfaces.CommentRepository, postRepository interfaces.PostRepository) *CommentService {
	return &CommentService{
		repository:     repository,
		postRepository: postRepository,
	}
}

func (this *CommentService) PublishOnPost(c *commands.PublishPostCommentCommand, a *auth.Auth) (*comments.Comment, *errors.Error) {
	if post := this.postRepository.FindById(c.PostId); post == nil {
		return nil, commenting.NewPostNotFoundError()
	}

	comment := comments.NewBuilder().
		Id(uuid.New()).
		AuthorId(a.UserId).
		Message(c.Message).
		PostId(c.PostId).
		CreatedAt(time.Now()).
		RaiseCreatedEvent().
		Build()

	err := this.repository.Save(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (this *CommentService) Remove(c *commands.RemoveCommentCommand, a *auth.Auth) (*comments.Comment, *errors.Error) {
	comment := this.repository.FindById(c.CommentId, true)
	if comment == nil {
		return nil, commenting.NewError("Такого комментария не существует")
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
