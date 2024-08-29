package domainfriendship

type DeniedEvent struct {
}

func (event *DeniedEvent) Type() EventType {
	return DENIED
}

func NewDeniedEvent() *DeniedEvent {
	return &DeniedEvent{}
}
