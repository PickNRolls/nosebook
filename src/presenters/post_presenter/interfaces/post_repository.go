package interfaces

import (
	"nosebook/src/presenters/post_presenter/dto"

	"github.com/google/uuid"
)

type PostRepository interface {
	FindAuthors(authorIds []uuid.UUID) ([]*dto.UserDTO, error)
	FindOwners(ownerIds []uuid.UUID) ([]*dto.UserDTO, error)
	FindLikers(likerIds []uuid.UUID) ([]*dto.UserDTO, error)
}
