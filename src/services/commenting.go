package services

import (
	"nosebook/src/domain/comments"
	"nosebook/src/errors"
	"nosebook/src/generics"
	"nosebook/src/services/auth"
	"nosebook/src/services/commenting/commands"
	"nosebook/src/services/commenting/interfaces"

	"github.com/google/uuid"
)

type CommentService struct {
	commentRepo interfaces.CommentRepository
}

func NewCommentService(commentRepo interfaces.CommentRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
	}
}

func (s *CommentService) FindByFilter(c *commands.FindCommentsCommand) *generics.SingleQueryResult[*comments.Comment] {
	return s.commentRepo.FindByFilter(c.Filter, c.Size)
}

func (s *CommentService) BatchFindByPostIds(ids []uuid.UUID) *generics.BatchQueryResult[*comments.Comment] {
	return s.commentRepo.FindByPostIds(ids)
}

func (s *CommentService) PublishOnPost(c *commands.PublishPostCommentCommand, a *auth.Auth) (*comments.Comment, *errors.Error) {
	comment, err := s.commentRepo.Create(
		comments.NewComment(a.UserId, c.Message).WithPost(c.PostId),
	)
	if err != nil {
		return nil, errors.From(err)
	}

	return comment, nil
}

func (s *CommentService) Remove(c *commands.RemoveCommentCommand, a *auth.Auth) (*comments.Comment, *errors.Error) {
	comment := s.commentRepo.FindById(c.CommentId, true)
	if comment == nil {
		return nil, errors.New("RemoveError", "Такого комментария не существует")
	}

	comment, error := comment.Remove(a.UserId)
	if error != nil {
		return nil, error
	}

	comment, err := s.commentRepo.Save(comment)
	if err != nil {
		return nil, errors.From(err)
	}

	return comment, nil
}

func (s *CommentService) Like(c *commands.LikeCommentCommand, a *auth.Auth) (*comments.Comment, *errors.Error) {
	comment := s.commentRepo.FindById(c.CommentId, true)
	if comment == nil {
		return nil, errors.New("LikeError", "Такого комментария не существует")
	}

	comment.Like(a.UserId)
	comment, err := s.commentRepo.Save(comment)
	if err != nil {
		return nil, errors.From(err)
	}

	return comment, nil
}
