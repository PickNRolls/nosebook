package dto

import (
	"time"

	"github.com/google/uuid"
)

type PostDTO struct {
	Id      uuid.UUID `json:"id"`
	Author  *UserDTO  `json:"author"`
	Owner   *UserDTO  `json:"owner"`
	Message string    `json:"message"`

	RecentComments QueryResultDTO[*CommentDTO] `json:"recentComments"`

	LikesCount       int        `json:"likesCount"`
	LikedByUser      bool       `json:"likedByUser"`
	RandomFiveLikers []*UserDTO `json:"randomFiveLikers"`

	CanBeRemovedByUser bool      `json:"canBeRemovedByUser"`
	CreatedAt          time.Time `json:"createdAt"`
}
