package permissionscomment

import "github.com/google/uuid"

func CanUpdateBy(comment Comment, userId uuid.UUID) *Error {
	if comment.AuthorId() != userId {
		return newError("Только автор комментария может его редактировать")
	}

	return nil
}

func CanRemoveBy(comment Comment, userId uuid.UUID) *Error {
	if comment.AuthorId() != userId {
		return newError("Только автор комментария может его удалить")
	}

	return nil
}
