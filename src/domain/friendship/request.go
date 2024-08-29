package domainfriendship

import (
	"nosebook/src/errors"
	"time"

	"github.com/google/uuid"
)

type FriendRequest struct {
	RequesterId uuid.UUID `json:"requesterId" db:"requester_id"`
	ResponderId uuid.UUID `json:"responderId" db:"responder_id"`
	Message     string    `json:"message" db:"message"`
	Accepted    bool      `json:"accepted" db:"accepted"`
	Viewed      bool      `json:"viewed" db:"viewed"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`

	events []Event
}

func New(
	requesterId uuid.UUID,
	responderId uuid.UUID,
	message string,
	accepted bool,
	viewed bool,
	createdAt time.Time,
	raiseCreatedEvent bool,
) *FriendRequest {
	output := &FriendRequest{
		RequesterId: requesterId,
		ResponderId: responderId,
		Message:     message,
		Accepted:    accepted,
		Viewed:      viewed,
		CreatedAt:   createdAt,
		events:      make([]Event, 0),
	}

	if raiseCreatedEvent {
		output.raiseEvent(NewCreatedEvent())
	}

	return output
}

func (this *FriendRequest) raiseEvent(event Event) {
	this.events = append(this.events, event)
}

func (this *FriendRequest) Events() []Event {
	return this.events
}

func (this *FriendRequest) ViewBy(userId uuid.UUID) *errors.Error {
	if userId != this.ResponderId {
		return errors.New("Friendship Error", "Только получатель может отметить заявку просмотренной")
	}

	this.Viewed = true
	this.raiseEvent(NewViewedEvent())

	return nil
}

func (this *FriendRequest) AcceptBy(userId uuid.UUID) *errors.Error {
	if userId != this.ResponderId {
		return errors.New("Friendship Error", "Только получатель может принять заявку")
	}

	this.Accepted = true
	this.Viewed = true
	this.raiseEvent(NewAcceptedEvent())

	return nil
}

func (this *FriendRequest) DenyBy(userId uuid.UUID) *errors.Error {
	if userId != this.ResponderId {
		return errors.New("Friendship Error", "Только получатель может отклонить заявку")
	}

	this.Accepted = false
	this.Viewed = true
	this.raiseEvent(NewDeniedEvent())

	return nil
}

func (this *FriendRequest) RemoveBy(userId uuid.UUID) *errors.Error {
	if !this.Accepted {
		return errors.New("Friendship Error", "Нельзя удалить непринятую заявку")
	}

	this.Accepted = false
	previousRequesterId := this.RequesterId
	previousResponderId := this.ResponderId

	if this.RequesterId == userId {
		this.RequesterId, this.ResponderId = this.ResponderId, this.RequesterId
	}

	this.raiseEvent(NewRemovedEvent(previousRequesterId, previousResponderId))

	return nil
}
