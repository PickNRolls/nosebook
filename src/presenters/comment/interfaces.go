package presentercomment

import (
	"nosebook/src/errors"
	"nosebook/src/services/auth"

	"github.com/google/uuid"
)

type userPresenter interface {
	FindByIds(ids uuid.UUIDs) ([]*user, *errors.Error)
}

type likePresenter interface {
	FindByPostIds(ids uuid.UUIDs, auth *auth.Auth) (map[uuid.UUID]*likes, *errors.Error)
}
