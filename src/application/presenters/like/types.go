package presenterlike

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type likesMap = map[uuid.UUID]*presenterdto.Likes
type usersMap = map[uuid.UUID]*presenterdto.User

type UserPresenter interface {
	FindByIds(ctx context.Context, ids uuid.UUIDs) (map[uuid.UUID]*presenterdto.User, *errors.Error)
}

type Resource interface {
	IDColumn() string
	Table() string
}

type dest struct {
	ResourceId uuid.UUID `db:"resource_id"`
	UserId     uuid.UUID `db:"user_id"`
}
