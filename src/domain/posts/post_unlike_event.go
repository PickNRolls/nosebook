package posts

import "github.com/google/uuid"

type PostUnlikeEvent struct {
	UserId uuid.UUID
}

func (event *PostUnlikeEvent) Type() PostEventType {
	return "UNLIKED"
}

func NewPostUnlikeEvent(userId uuid.UUID) *PostUnlikeEvent {
	return &PostUnlikeEvent{
		UserId: userId,
	}
}
