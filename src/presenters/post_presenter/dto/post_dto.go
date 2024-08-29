package dto

import (
	"time"

	"nosebook/src/domain/comment"
	"nosebook/src/presenters/dto"

	"github.com/google/uuid"
)

type PostDTO struct {
	Id      uuid.UUID `json:"id"`
	Author  *UserDTO  `json:"author"`
	Owner   *UserDTO  `json:"owner"`
	Message string    `json:"message"`

	RecentComments dto.SingleQueryResultDTO[*domaincomment.Comment] `json:"recentComments"`

	LikesCount       int        `json:"likesCount"`
	LikedByUser      bool       `json:"likedByUser"`
	RandomFiveLikers []*UserDTO `json:"randomFiveLikers"`

	CanBeRemovedByUser bool      `json:"canBeRemovedByUser"`
	CreatedAt          time.Time `json:"createdAt"`
}
