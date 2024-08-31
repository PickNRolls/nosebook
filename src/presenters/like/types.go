package presenterlike

import (
	"nosebook/src/errors"
	presenterdto "nosebook/src/presenters/dto"

	"github.com/google/uuid"
)

type likesMap = map[uuid.UUID]*presenterdto.Likes
type usersMap = map[uuid.UUID]*presenterdto.User

type UserPresenter interface {
	FindByIds(ids uuid.UUIDs) ([]*presenterdto.User, *errors.Error)
}

type Resource interface {
	IDColumn() string
	Table() string
}

type dest struct {
	ResourceId uuid.UUID `db:"resource_id"`
	UserId     uuid.UUID `db:"user_id"`
}
