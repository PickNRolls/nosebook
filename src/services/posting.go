package services

import (
	"errors"
	"nosebook/src/domain/posts"
	"nosebook/src/services/auth"
	"nosebook/src/services/posting/commands"
	"nosebook/src/services/posting/interfaces"
)

type PostingService struct {
	postsRepo interfaces.PostsRepository
}

func NewPostingService(postsRepo interfaces.PostsRepository) *PostingService {
	return &PostingService{
		postsRepo: postsRepo,
	}
}

func (s *PostingService) Publish(c *commands.PublishPostCommand, a *auth.Auth) (*posts.Post, error) {
	post := posts.NewPost(a.UserId, c.OwnerId, c.Message)
	return s.postsRepo.Create(post)
}

func (s *PostingService) Remove(c *commands.RemovePostCommand, a *auth.Auth) (*posts.Post, error) {
	post := s.postsRepo.FindById(c.Id)
	if post == nil {
		return nil, errors.New("No such post.")
	}

	if post.OwnerId != a.UserId {
		return nil, errors.New("You can't remove this post, only post owners can remove posts.")
	}

	_, err := s.postsRepo.Remove(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}
