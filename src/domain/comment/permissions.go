package domaincomment

import "github.com/google/uuid"

type permissions interface {
	CanRemoveBy(comment *Comment, userId uuid.UUID) *Error
	CanUpdateBy(comment *Comment, userId uuid.UUID) *Error
}

type defaultPermissions struct{}

func (this *defaultPermissions) CanUpdateBy(comment *Comment, userId uuid.UUID) *Error {
	if comment.AuthorId != userId {
		return NewError("Только автор комментария может его редактировать")
	}

	return nil
}

func (this *defaultPermissions) CanRemoveBy(comment *Comment, userId uuid.UUID) *Error {
	if comment.AuthorId != userId {
		return NewError("Только автор комментария может его удалить")
	}

	return nil
}
