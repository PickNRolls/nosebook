package like

import (
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
)

type Service struct {
	repository         Repository
	notifierRepository NotifierRepository
}

func New(repository Repository, notifierRepository NotifierRepository) *Service {
	return &Service{
		repository:         repository,
		notifierRepository: notifierRepository,
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

	err = this.repository.Save(like)
	if err != nil {
		return nil, err
	}

	if like.Value {
		notifier := this.notifierRepository.FindByUserId(like.Resource.Owner().Id())
		if notifier != nil {
			err = notifier.NotifyAbout(like)
			if err != nil {
				return nil, err
			}
		}
	}

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

	err = this.repository.Save(like)
	if err != nil {
		return nil, err
	}

	if like.Value {
		notifier := this.notifierRepository.FindByUserId(like.Resource.Owner().Id())
		if notifier != nil {
			err = notifier.NotifyAbout(like)
			if err != nil {
				return nil, err
			}
		}
	}

	return &resultData{
		CommentId: &c.Id,
		Liked:     like.Value,
	}, nil
}
