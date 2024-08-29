package domainfriendship

type ViewedEvent struct {
}

func (event *ViewedEvent) Type() EventType {
	return VIEWED
}

func NewViewedEvent() *ViewedEvent {
	return &ViewedEvent{}
}
