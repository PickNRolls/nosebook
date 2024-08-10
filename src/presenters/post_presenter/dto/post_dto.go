package dto

import (
	"time"

	"github.com/google/uuid"
)

type PostDTO struct {
	Id                 uuid.UUID  `json:"id"`
	Author             *UserDTO   `json:"author"`
	Owner              *UserDTO   `json:"owner"`
	Message            string     `json:"message"`
	CreatedAt          time.Time  `json:"createdAt"`
	LikesCount         int        `json:"likesCount"`
	LikedByUser        bool       `json:"likedByUser"`
	CanBeRemovedByUser bool       `json:"canBeRemovedByUser"`
	RandomFiveLikers   []*UserDTO `json:"randomFiveLikers"`
}
