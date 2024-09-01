package domainpost

import (
	"database/sql"
	"nosebook/src/lib/clock"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id        uuid.UUID
	AuthorId  uuid.UUID
	OwnerId   uuid.UUID
	Message   string
	CreatedAt time.Time
	RemovedAt sql.NullTime

	events []PostEvent
}

func NewPost(
	id uuid.UUID,
	authorId uuid.UUID,
	ownerId uuid.UUID,
	message string,
	createdAt time.Time,
	removedAt sql.NullTime,
	raiseCreatedEvent bool,
) *Post {
	post := &Post{
		Id:        id,
		AuthorId:  authorId,
		OwnerId:   ownerId,
		Message:   message,
		CreatedAt: createdAt,
		RemovedAt: removedAt,

		events: make([]PostEvent, 0),
	}
	if raiseCreatedEvent {
		post.raiseEvent(NewPostCreatedEvent())
	}

	return post
}

func (post *Post) Events() []PostEvent {
	return post.events
}

func (post *Post) raiseEvent(event PostEvent) {
	post.events = append(post.events, event)
}

func (post *Post) CanBeRemovedBy(userId uuid.UUID) *PostError {
	if post.OwnerId == userId || post.AuthorId == userId {
		return nil
	}

	return NewError("Только автор и владелец поста может его удалить")
}

func (post *Post) RemoveBy(userId uuid.UUID) *PostError {
	err := post.CanBeRemovedBy(userId)
	if err != nil {
		return err
	}

	if post.RemovedAt.Valid {
		return NewError("Пост уже удален")
	}

	post.RemovedAt = sql.NullTime{
		Time:  clock.Now(),
		Valid: true,
	}
	post.raiseEvent(NewPostRemovedEvent())

	return nil
}

func (post *Post) CanBeEditedBy(userId uuid.UUID) *PostError {
	if post.AuthorId == userId {
		return nil
	}

	return NewError("Только владелец поста может его редактировать")
}

func (post *Post) EditBy(userId uuid.UUID, message string) *PostError {
	err := post.CanBeEditedBy(userId)
	if err != nil {
		return err
	}

	post.Message = message
	post.raiseEvent(NewPostEditedEvent(post.Message))

	return nil
}
