package conversation

import (
	"github.com/google/uuid"
)

type SendMessageCommand struct {
	RecipientId uuid.UUID `json:"recipientId"`
	Text        string    `json:"text"`
	ReplyTo     uuid.UUID `json:"replyTo"`
}
