package presenterfriendship

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type user = presenterdto.User

type UserPresenter interface {
	FindByIds(ids uuid.UUIDs) ([]*user, *errors.Error)
}
