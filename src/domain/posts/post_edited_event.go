package posts

type PostEditedEvent struct {
	Message string
}

func (event *PostEditedEvent) Type() PostEventType {
	return EDITED
}

func NewPostEditedEvent(message string) *PostEditedEvent {
	return &PostEditedEvent{
		Message: message,
	}
}
