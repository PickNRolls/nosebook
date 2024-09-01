package presentercomment

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type LikePresenter interface {
	FindByCommentIds(ids uuid.UUIDs, auth *auth.Auth) (map[uuid.UUID]*likes, *errors.Error)
}

type UserPresenter interface {
	FindByIds(ids uuid.UUIDs) ([]*presenterdto.User, *errors.Error)
}

type Permissions interface {
	CanRemoveBy(comment *Dest, userId uuid.UUID) bool
	CanUpdateBy(comment *Dest, userId uuid.UUID) bool
}

type user = presenterdto.User
type likes = presenterdto.Likes
type comment = presenterdto.Comment

type FindByFilterInput struct {
	Ids    []string
	PostId string

	Next  string
	Prev  string
	Limit uint64
	Last  bool
}

type FindByFilterOutput = presenterdto.FindOut[*comment]
