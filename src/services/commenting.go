package services

import (
	"errors"
	"nosebook/src/domain/comments"
	"nosebook/src/services/auth"
	"nosebook/src/services/commenting/commands"
	"nosebook/src/services/commenting/interfaces"
)

type CommentService struct {
	commentRepo interfaces.CommentRepository
}

func NewCommentService(commentRepo interfaces.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

func (s *CommentService) PublishOnPost(c *commands.PublishPostCommentCommand, a *auth.Auth) (*comments.Comment, error) {
	comment, err := s.commentRepo.Create(
		comments.NewComment(a.UserId, c.Message).WithPost(c.PostId),
	)
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

	comment.Like(a.UserId)
	comment, err := s.commentRepo.Save(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
