package presenterlike

import (
	"nosebook/src/errors"
	presenterdto "nosebook/src/presenters/dto"

	"github.com/google/uuid"
)

type userPresenter interface {
	FindByIds(ids uuid.UUIDs) ([]*presenterdto.User, *errors.Error)
}
