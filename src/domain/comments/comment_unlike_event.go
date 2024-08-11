package comments

import "github.com/google/uuid"

type CommentUnlikeEvent struct {
	UserId uuid.UUID
}

func (event *CommentUnlikeEvent) Type() CommentEventType {
	return UNLIKED
}

func NewCommentUnlikeEvent(userId uuid.UUID) *CommentUnlikeEvent {
	return &CommentUnlikeEvent{
		UserId: userId,
	}
}
