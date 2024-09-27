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
	AvatarUrl      string    `json:"avatarUrl,omitempty" db:"avatar_url"`
	LastActivityAt time.Time `json:"lastActivityAt" db:"last_activity_at"`
}

func New(firstName string, lastName string, nick string, passhash string, avatarUrl string) *User {
	return &User{
		Id:             uuid.New(),
		FirstName:      firstName,
		LastName:       lastName,
		Nick:           nick,
		Passhash:       passhash,
		CreatedAt:      clock.Now(),
		AvatarUrl:      avatarUrl,
		LastActivityAt: clock.Now(),
	}
}

func (this *User) MarkActivity() {
	this.LastActivityAt = clock.Now()
}

func (this *User) ChangeAvatar(url string) {
	this.AvatarUrl = url
}
