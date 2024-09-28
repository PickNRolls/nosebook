package presenterdto

import (
	"time"

	"github.com/google/uuid"
)

type UserAvatar struct {
	Url       string    `json:"url"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	Id        uuid.UUID `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Nick      string    `json:"nick" db:"nick"`

	LastOnlineAt time.Time   `json:"lastOnlineAt"`
	Online       bool        `json:"online"`
	Avatar       *UserAvatar `json:"avatar,omitempty"`
}
