package services

import (
	"errors"
	"nosebook/src/domain/comments"
	"nosebook/src/services/auth"
	"nosebook/src/services/commenting/commands"
	"nosebook/src/services/commenting/interfaces"
)

type CommentService struct {
	commentRepo      interfaces.CommentRepository
	commentLikesRepo interfaces.CommentLikesRepository
}

func NewCommentService(commentRepo interfaces.CommentRepository, commentLikesRepo interfaces.CommentLikesRepository) *CommentService {
	return &CommentService{
		commentRepo:      commentRepo,
		commentLikesRepo: commentLikesRepo,
	}
}

func (s *CommentService) PublishOnPost(c *commands.PublishPostCommentCommand, a *auth.Auth) (*comments.Comment, error) {
	comment := comments.NewComment(a.UserId, c.Message)
	comment, err := s.commentRepo.CreateForPost(c.PostId, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) Remove(c *commands.RemoveCommentCommand, a *auth.Auth) (*comments.Comment, error) {
	comment, err := s.commentRepo.Remove(c.CommentId)
	if err != nil {
		return nil, err
	}

	return comment, err
}

func (s *CommentService) Like(c *commands.LikeCommentCommand, a *auth.Auth) (*comments.Comment, error) {
	comment := s.commentRepo.FindById(c.CommentId)
	if comment == nil {
		return nil, errors.New("No such comment.")
	}

	like := s.commentLikesRepo.Find(c.CommentId, a.UserId)
	if like == nil {
		_, err := s.commentLikesRepo.Create(comments.NewCommentLike(c.CommentId, a.UserId))
		if err != nil {
			return nil, err
		}
	} else {
		_, err := s.commentLikesRepo.Remove(like)
		if err != nil {
			return nil, err
		}
	}

	return comment, nil
}
