package comments

import "github.com/google/uuid"

type CommentLikeEvent struct {
	UserId uuid.UUID
}

func (event *CommentLikeEvent) Type() CommentEventType {
	return LIKED
}

func NewCommentLikeEvent(userId uuid.UUID) *CommentLikeEvent {
	return &CommentLikeEvent{
		UserId: userId,
	}
}
