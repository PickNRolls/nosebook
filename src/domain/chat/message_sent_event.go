package domainchat

type MessageSentEvent struct {
	Message *message
}

func (event *MessageSentEvent) Type() EventType {
	return MESSAGE_SENT
}

func newMessageSentEvent(message *message) *MessageSentEvent {
	return &MessageSentEvent{
		Message: message,
	}
}
