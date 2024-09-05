package domainmessage

import "github.com/google/uuid"

type Permissions interface {
	CanRemoveBy(message *Message, userId uuid.UUID) *Error
	CanUpdateBy(message *Message, userId uuid.UUID) *Error
}

type defaultPermissions struct{}

func (this *defaultPermissions) CanUpdateBy(message *Message, userId uuid.UUID) *Error {
	if message.AuthorId != userId {
		return newError("Только автор сообщения может его редактировать")
	}

	return nil
}

func (this *defaultPermissions) CanRemoveBy(message *Message, userId uuid.UUID) *Error {
	if message.AuthorId != userId {
		return newError("Только автор сообщения может его удалить")
	}

	return nil
}
