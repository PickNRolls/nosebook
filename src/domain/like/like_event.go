package domainlike

type LikeEvent struct{}

func NewLikeEvent() *LikeEvent {
	return &LikeEvent{}
}

func (this *LikeEvent) Type() EventType {
	return LIKED
}
