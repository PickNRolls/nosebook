package like

import (
	"fmt"
	domainlike "nosebook/src/domain/like"
	"nosebook/src/services/auth"
)

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (this *Service) LikePost(c *LikePostCommand, a *auth.Auth) (*domainlike.Like, *Error) {
	like, err := this.repository.
		WithPostId(c.Id).
		WithUserId(a.UserId).
		FindOne()
	if err != nil {
		return nil, err
	}

	fmt.Println(like)

	err = like.Toggle()
	if err != nil {
		return nil, err
	}

	fmt.Println(like)

	this.repository.Save(like)

	return like, nil
}

func (this *Service) LikeComment(c *LikeCommentCommand, a *auth.Auth) (*domainlike.Like, *Error) {
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

	return like, nil
}
