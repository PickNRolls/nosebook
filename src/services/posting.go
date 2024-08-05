package services

import (
	// "errors"
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
