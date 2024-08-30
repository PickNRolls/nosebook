package presenterdto

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Nick      string    `json:"nick" db:"nick"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
