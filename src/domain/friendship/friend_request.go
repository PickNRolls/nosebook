package friendship

import (
	"github.com/google/uuid"
	"time"
)

type FriendRequest struct {
	RequesterId uuid.UUID `json:"requesterId" db:"requester_id"`
	ResponderId uuid.UUID `json:"responderId" db:"responder_id"`
	Message     string    `json:"message" db:"message"`
	Accepted    bool      `json:"accepted" db:"accepted"`
	Viewed      bool      `json:"viewed" db:"viewed"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

func NewFriendRequest(requesterId uuid.UUID, responderId uuid.UUID, message string) *FriendRequest {
	return &FriendRequest{
		RequesterId: requesterId,
		ResponderId: responderId,
		Message:     message,
		Accepted:    false,
		Viewed:      false,
		CreatedAt:   time.Now(),
	}
}
