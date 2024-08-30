package presenterdto

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Author    *User     `json:"author"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
