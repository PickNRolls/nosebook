package presenterfriendship

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type user = presenterdto.User

type FindByFilterInput struct {
	UserId     string
	Text       string
	OnlyMutual bool
	OnlyOnline bool
	Shuffle    bool

	Next  string
	Prev  string
	Limit uint64
}

type FindByFilterOutput = presenterdto.FindOut[*user]

type UserPresenter interface {
	FindByIds(ids uuid.UUIDs) ([]*user, *errors.Error)
}
