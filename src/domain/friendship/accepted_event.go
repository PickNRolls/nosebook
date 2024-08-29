package domainfriendship

type AcceptedEvent struct {
}

func (event *AcceptedEvent) Type() EventType {
	return ACCEPTED
}

func NewAcceptedEvent() *AcceptedEvent {
	return &AcceptedEvent{}
}
