package posts

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id        uuid.UUID   `json:"id" db:"id"`
	AuthorId  uuid.UUID   `json:"authorId" db:"author_id"`
	OwnerId   uuid.UUID   `json:"ownerId" db:"owner_id"`
	Message   string      `json:"message" db:"message"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`
	RemovedAt time.Time   `json:"removedAt" db:"removed_at"`
	LikedBy   []uuid.UUID `json:"-"`

	Events []PostEvent `json:"-"`
}

func NewPost(authorId uuid.UUID, ownerId uuid.UUID, message string) *Post {
	return &Post{
		Id:        uuid.New(),
		AuthorId:  authorId,
		OwnerId:   ownerId,
		Message:   message,
		CreatedAt: time.Now(),
		RemovedAt: time.Time{},
		LikedBy:   make([]uuid.UUID, 0),

		Events: make([]PostEvent, 0),
	}
}

func (post *Post) Like(userId uuid.UUID) *Post {
	for i, id := range post.LikedBy {
		if id == userId {
			post.LikedBy[i] = post.LikedBy[len(post.LikedBy)-1]
			post.LikedBy = post.LikedBy[:len(post.LikedBy)-1]
			post.Events = append(post.Events, NewPostUnlikeEvent(userId))
			return post
		}
	}

	post.LikedBy = append(post.LikedBy, userId)
	post.Events = append(post.Events, NewPostLikeEvent(userId))
	return post
}
