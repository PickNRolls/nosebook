package services

import (
	"errors"
	"nosebook/src/domain/posts"
	"nosebook/src/services/auth"
	"nosebook/src/services/posting/commands"
	"nosebook/src/services/posting/interfaces"
)

type PostingService struct {
	postRepo      interfaces.PostRepository
	postLikesRepo interfaces.PostLikesRepository
}

func NewPostingService(postRepo interfaces.PostRepository, postLikesRepo interfaces.PostLikesRepository) *PostingService {
	return &PostingService{
		postRepo:      postRepo,
		postLikesRepo: postLikesRepo,
	}
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

	like := s.postLikesRepo.Find(c.PostId, a.UserId)
	if like == nil {
		_, err := s.postLikesRepo.Create(posts.NewPostLike(c.PostId, a.UserId))
		if err != nil {
			return nil, err
		}
	} else {
		_, err := s.postLikesRepo.Remove(like)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}
