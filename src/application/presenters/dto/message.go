package presenterdto

import (
	"nosebook/src/lib/nullable"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id        uuid.UUID      `json:"id"`
	Author    *User          `json:"author"`
	Text      string         `json:"text"`
	ReplyTo   *nullable.UUID `json:"replyTo,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}
