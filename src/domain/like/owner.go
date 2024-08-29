package domainlike

import "github.com/google/uuid"

type OwnerType string

const (
	USER_OWNER OwnerType = "user"
)

type Owner interface {
	Id() uuid.UUID
	Type() OwnerType
}

type UserOwner struct {
	id uuid.UUID
}

func NewUserOwner(id uuid.UUID) *UserOwner {
	return &UserOwner{
		id: id,
	}
}

func (this *UserOwner) Id() uuid.UUID {
	return this.id
}

func (this *UserOwner) Type() OwnerType {
	return USER_OWNER
}
