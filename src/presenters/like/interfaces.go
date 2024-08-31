package presenterlike

import (
	"nosebook/src/errors"
	presenterdto "nosebook/src/presenters/dto"

	"github.com/google/uuid"
)

type UserPresenter interface {
	FindByIds(ids uuid.UUIDs) ([]*presenterdto.User, *errors.Error)
}

type Resource interface {
	IDColumn() string
	Table() string
}
