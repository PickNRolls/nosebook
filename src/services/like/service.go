package like

import (
	commandresult "nosebook/src/lib/command_result"
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

func (this *Service) LikePost(c *LikePostCommand, a *auth.Auth) *commandresult.Result {
	like, err := this.repository.
		WithPostId(c.Id).
		WithUserId(a.UserId).
		FindOne()
	if err != nil {
		return commandresult.Fail(err)
	}

	err = like.Toggle()
	if err != nil {
		return commandresult.Fail(err)
	}

	this.repository.Save(like)

	return commandresult.Ok().WithData(resultData{
		PostId: &c.Id,
		Liked:  like.Value,
	})
}

func (this *Service) LikeComment(c *LikeCommentCommand, a *auth.Auth) *commandresult.Result {
	like, err := this.repository.
		WithCommentId(c.Id).
		WithUserId(a.UserId).
		FindOne()
	if err != nil {
		return commandresult.Fail(err)
	}

	err = like.Toggle()
	if err != nil {
		return commandresult.Fail(err)
	}

	this.repository.Save(like)

	return commandresult.Ok().WithData(resultData{
		CommentId: &c.Id,
		Liked:     like.Value,
	})
}
