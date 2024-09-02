package permissionspost

import "github.com/google/uuid"

func CanUpdateBy(post Post, userId uuid.UUID) *Error {
	if post.AuthorId() != userId {
		return newError("Только автор поста может его редактировать")
	}

	return nil
}

func CanRemoveBy(post Post, userId uuid.UUID) *Error {
	if post.AuthorId() != userId && post.OwnerId() != userId {
		return newError("Только автор или владелец поста может его удалить")
	}

	return nil
}
