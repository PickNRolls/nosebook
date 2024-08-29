package presenterpost

import (
	"nosebook/src/errors"
	"time"

	"github.com/google/uuid"
)

type FindByFilterInput struct {
	OwnerId  string
	AuthorId string
	Cursor   string
}

type FindByFilterOutput struct {
	Err   *errors.Error `json:"error,omitempty"`
	Posts []*PostDTO    `json:"data,omitempty"`
	Next  string        `json:"next,omitempty"`
}

type PostDTO struct {
	Id        uuid.UUID `json:"id"`
	Author    *UserDTO  `json:"author"`
	Owner     *UserDTO  `json:"owner"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserDTO struct {
	Id        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Nick      string    `json:"nick" db:"nick"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
