package conversation

import (
	"nosebook/src/lib/nullable"

	"github.com/google/uuid"
)

type SendMessageCommand struct {
	RecipientId uuid.UUID     `json:"recipientId"`
	Text        string        `json:"text"`
	ReplyTo     nullable.UUID `json:"replyTo"`
}
