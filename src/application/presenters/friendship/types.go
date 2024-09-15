package presenterfriendship

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type user = presenterdto.User

type RequestType = string

const (
	INCOMING  RequestType = "incoming"
	OUTCOMING RequestType = "outcoming"
)

type Request struct {
	Type     RequestType `json:"type"`
	Accepted bool        `json:"accepted"`
	User     *user       `json:"user"`
}

type UserPresenter interface {
	FindByIds(ctx context.Context, ids uuid.UUIDs) (map[uuid.UUID]*user, *errors.Error)
}
