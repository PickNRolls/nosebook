package presentercomment

import (
	presenterdto "nosebook/src/presenters/dto"
	"time"

	"github.com/google/uuid"
)

type FindByFilterInput struct {
	PostId string
	Next   string
	Prev   string
	Last   bool
}

type FindByFilterOutput = presenterdto.FindOut[*presenterdto.Comment]

type userDTO struct {
	Id        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Nick      string    `json:"nick" db:"nick"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
