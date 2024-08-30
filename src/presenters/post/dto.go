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
	Posts []*postDTO    `json:"data,omitempty"`
	Next  string        `json:"next,omitempty"`
}

type postDTO struct {
	Id        uuid.UUID `json:"id"`
	Author    *userDTO  `json:"author"`
	Owner     *userDTO  `json:"owner"`
	Message   string    `json:"message"`
	Likes     *likesDTO `json:"likes"`
	CreatedAt time.Time `json:"createdAt"`
}

type likesDTO struct {
	Count            int        `json:"count"`
	RandomFiveLikers []*userDTO `json:"randomFiveLikers"`
	Liked            bool       `json:"liked"`
}

type userDTO struct {
	Id        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Nick      string    `json:"nick" db:"nick"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
