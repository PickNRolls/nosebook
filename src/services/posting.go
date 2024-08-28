package services

import (
	"nosebook/src/domain/posts"
	"nosebook/src/errors"
	"nosebook/src/services/auth"
	"nosebook/src/services/posting"
	"nosebook/src/services/posting/commands"
	"nosebook/src/services/posting/interfaces"
	"time"

	"github.com/google/uuid"
)

type PostingService struct {
	repository interfaces.PostRepository
}

func NewPostingService(repository interfaces.PostRepository) *PostingService {
	return &PostingService{
		repository: repository,
	}
}

func (s *PostingService) Publish(c *commands.PublishPostCommand, a *auth.Auth) (*posts.Post, *errors.Error) {
	post := posts.NewBuilder().
		Id(uuid.New()).
		AuthorId(a.UserId).
		OwnerId(c.OwnerId).
		Message(c.Message).
		CreatedAt(time.Now()).
		RaiseCreatedEvent().
		Build()

	err := s.repository.Save(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostingService) Remove(c *commands.RemovePostCommand, a *auth.Auth) (*posts.Post, *errors.Error) {
	post := s.repository.FindById(c.Id)
	if post == nil {
		return nil, posting.NewNotFoundError()
	}

	err := post.RemoveBy(a.UserId)
	if err != nil {
		return nil, err
	}

	err = s.repository.Save(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostingService) Edit(c *commands.EditPostCommand, a *auth.Auth) (*posts.Post, *errors.Error) {
	post := s.repository.FindById(c.Id)
	if post == nil {
		return nil, posting.NewNotFoundError()
	}

	err := post.EditBy(a.UserId, c.Message)
	if err != nil {
		return nil, err
	}

	err = s.repository.Save(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}
