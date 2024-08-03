package users

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Nick      string    `json:"nick" db:"nick"`
	Passhash  string    `json:"passhash" db:"passhash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func NewUser(firstName string, lastName string, nick string, passhash string) *User {
	return &User{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Nick:      nick,
		Passhash:  passhash,
		CreatedAt: time.Now(),
	}
}
