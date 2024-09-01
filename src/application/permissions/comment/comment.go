package permissionscomment

import "github.com/google/uuid"

func CanUpdateBy(comment CommentToUpdate, userId uuid.UUID) *Error {
	if comment.AuthorId() != userId {
		return newError("Вы не можете редактировать комментарий")
	}

	return nil
}

func CanRemoveBy(comment CommentToRemove, userId uuid.UUID) *Error {
	if comment.AuthorId() != userId && comment.ResourceOwnerId() != userId {
		return newError("Вы не можете удалить комментарий")
	}

	return nil
}
