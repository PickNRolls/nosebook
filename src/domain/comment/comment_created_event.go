package domaincomment

type CommentCreatedEvent struct {
}

func (event *CommentCreatedEvent) Type() CommentEventType {
	return CREATED
}

func NewCommentCreatedEvent() *CommentCreatedEvent {
	return &CommentCreatedEvent{}
}
