package services

import (
	"errors"
	"nosebook/src/domain/posts"
	"nosebook/src/services/auth"
	"nosebook/src/services/posting/commands"
	"nosebook/src/services/posting/interfaces"
	"nosebook/src/services/posting/structs"
)

type PostingService struct {
	postRepo interfaces.PostRepository
}

func NewPostingService(postRepo interfaces.PostRepository) *PostingService {
	return &PostingService{
		postRepo: postRepo,
	}
}

func (s *PostingService) FindByFilter(c *commands.FindPostsCommand, a *auth.Auth) structs.QueryResult {
	return s.postRepo.FindByFilter(c.Filter)
}

func (s *PostingService) Publish(c *commands.PublishPostCommand, a *auth.Auth) (*posts.Post, error) {
	post := posts.NewPost(a.UserId, c.OwnerId, c.Message)
	return s.postRepo.Create(post)
}

func (s *PostingService) Remove(c *commands.RemovePostCommand, a *auth.Auth) (*posts.Post, error) {
	post := s.postRepo.FindById(c.Id)
	if post == nil {
		return nil, errors.New("No such post.")
	}

	if post.OwnerId != a.UserId {
		return nil, errors.New("You can't remove this post, only post owners can remove posts.")
	}

	_, err := s.postRepo.Remove(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostingService) Like(c *commands.LikePostCommand, a *auth.Auth) (*posts.Post, error) {
	post := s.postRepo.FindById(c.PostId)
	if post == nil {
		return nil, errors.New("No such post.")
	}

	post.Like(a.UserId)
	post, err := s.postRepo.Save(post)

	if err != nil {
		return nil, err
	}

	return post, nil
}
