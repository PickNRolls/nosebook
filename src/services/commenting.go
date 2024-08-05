package services

import (
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
	comment := comments.NewComment(a.UserId, c.Message)
	comment, err := s.commentRepo.CreateForPost(c.PostId, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
