package like

import "nosebook/src/errors"

type Error = errors.Error

func NewError(message string) *Error {
	return errors.New("LikeError", message)
}

func NewPostNotFoundError() *Error {
	return NewError("Такого поста не существует")
}

func NewCommentNotFoundError() *Error {
	return NewError("Такого комментария не существует")
}
