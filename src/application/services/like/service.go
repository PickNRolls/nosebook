package like

import (
	"context"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
)

type Service struct {
	repository Repository
	notifier   Notifier
}

func New(repository Repository, notifier Notifier) *Service {
	return &Service{
		repository: repository,
		notifier:   notifier,
	}
}

func (this *Service) LikePost(parent context.Context, c LikePostCommand, auth *auth.Auth) (*resultData, *errors.Error) {
	like, err := this.repository.
		WithPostId(c.Id).
		WithUserId(auth.UserId).
		FindOne()
	if err != nil {
		return nil, err
	}

	err = like.Toggle()
	if err != nil {
		return nil, err
	}

	err = this.repository.Save(like)
	if err != nil {
		return nil, err
	}

	if like.Value && like.Resource.Owner().Id() != auth.UserId {
		err = this.notifier.NotifyAbout(like.Resource.Owner().Id(), like)
		if err != nil {
			return nil, err
		}
	}

	return &resultData{
		PostId: &c.Id,
		Liked:  like.Value,
	}, nil
}

func (this *Service) LikeComment(parent context.Context, c LikeCommentCommand, auth *auth.Auth) (*resultData, *errors.Error) {
	like, err := this.repository.
		WithCommentId(c.Id).
		WithUserId(auth.UserId).
		FindOne()
	if err != nil {
		return nil, err
	}

	err = like.Toggle()
	if err != nil {
		return nil, err
	}

	err = this.repository.Save(like)
	if err != nil {
		return nil, err
	}

	if like.Value && like.Resource.Owner().Id() != auth.UserId {
		err = this.notifier.NotifyAbout(like.Resource.Owner().Id(), like)
		if err != nil {
			return nil, err
		}
	}

	return &resultData{
		CommentId: &c.Id,
		Liked:     like.Value,
	}, nil
}
