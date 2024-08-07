package posts

import "github.com/google/uuid"

type PostLikeEvent struct {
	UserId uuid.UUID
}

func (event *PostLikeEvent) Type() PostEventType {
	return "LIKED"
}

func NewPostLikeEvent(userId uuid.UUID) *PostLikeEvent {
	return &PostLikeEvent{
		UserId: userId,
	}
}
