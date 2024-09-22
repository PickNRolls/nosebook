package domainuser

import (
	"nosebook/src/lib/clock"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id             uuid.UUID `json:"id" db:"id"`
	FirstName      string    `json:"firstName" db:"first_name"`
	LastName       string    `json:"lastName" db:"last_name"`
	Nick           string    `json:"nick" db:"nick"`
	Passhash       string    `json:"passhash" db:"passhash"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	LastActivityAt time.Time `json:"lastActivityAt" db:"last_activity_at"`
}

func New(firstName string, lastName string, nick string, passhash string) *User {
	return &User{
		Id:             uuid.New(),
		FirstName:      firstName,
		LastName:       lastName,
		Nick:           nick,
		Passhash:       passhash,
		CreatedAt:      clock.Now(),
		LastActivityAt: clock.Now(),
	}
}

func (this *User) MarkActivity() {
  this.LastActivityAt = clock.Now()
}
