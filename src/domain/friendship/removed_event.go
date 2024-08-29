package domainfriendship

import "github.com/google/uuid"

type RemovedEvent struct {
	PreviousRequesterId uuid.UUID
	PreviousResponderId uuid.UUID
}

func (event *RemovedEvent) Type() EventType {
	return REMOVED
}

func NewRemovedEvent(previousRequesterId uuid.UUID, previousResponderId uuid.UUID) *RemovedEvent {
	return &RemovedEvent{
		PreviousRequesterId: previousRequesterId,
		PreviousResponderId: previousResponderId,
	}
}
