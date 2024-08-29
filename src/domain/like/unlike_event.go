package domainlike

type UnlikeEvent struct{}

func NewUnlikeEvent() *UnlikeEvent {
	return &UnlikeEvent{}
}

func (this *UnlikeEvent) Type() EventType {
	return UNLIKED
}
