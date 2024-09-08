package like

import (
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
)

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (this *Service) LikePost(c *LikePostCommand, a *auth.Auth) (*resultData, *errors.Error) {
	like, err := this.repository.
		WithPostId(c.Id).
		WithUserId(a.UserId).
		FindOne()
	if err != nil {
		return nil, err
	}

	err = like.Toggle()
	if err != nil {
		return nil, err
	}

	this.repository.Save(like)

	return &resultData{
		PostId: &c.Id,
		Liked:  like.Value,
	}, nil
}

func (this *Service) LikeComment(c *LikeCommentCommand, a *auth.Auth) (*resultData, *errors.Error) {
	like, err := this.repository.
		WithCommentId(c.Id).
		WithUserId(a.UserId).
		FindOne()
	if err != nil {
		return nil, err
	}

	err = like.Toggle()
	if err != nil {
		return nil, err
	}

	this.repository.Save(like)

	return &resultData{
		CommentId: &c.Id,
		Liked:     like.Value,
	}, nil
}
