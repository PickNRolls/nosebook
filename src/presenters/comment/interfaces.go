package presentercomment

import (
	"nosebook/src/errors"
	presenterdto "nosebook/src/presenters/dto"
	"nosebook/src/services/auth"

	"github.com/google/uuid"
)

type likePresenter interface {
	FindByCommentIds(ids uuid.UUIDs, auth *auth.Auth) (map[uuid.UUID]*likes, *errors.Error)
}

type userPresenter interface {
	FindByIds(ids uuid.UUIDs) ([]*presenterdto.User, *errors.Error)
}
