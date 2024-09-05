package presenterdto

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	Id           uuid.UUID `json:"id"`
	Interlocutor *User     `json:"interlocutor"`
	LastMessage  *Message  `json:"lastMessage"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (this *Conversation) IsChat() {}
