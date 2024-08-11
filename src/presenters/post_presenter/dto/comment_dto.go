package dto

import (
	"time"

	"github.com/google/uuid"
)

type CommentDTO struct {
	Id        uuid.UUID `json:"id"`
	Author    *UserDTO  `json:"author"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`

	LikesCount       int        `json:"likesCount"`
	LikedByUser      bool       `json:"likedByUser"`
	RandomFiveLikers []*UserDTO `json:"randomFiveLikers"`

	CanRemove bool `json:"canRemove"`
}
