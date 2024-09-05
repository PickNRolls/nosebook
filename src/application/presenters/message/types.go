package presentermessage

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type message = presenterdto.Message
type user = presenterdto.User

type UserPresenter interface {
	FindByIds(ids uuid.UUIDs) (map[uuid.UUID]*user, *errors.Error)
}
