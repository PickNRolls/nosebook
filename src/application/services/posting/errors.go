package posting

import "nosebook/src/errors"

type Error = errors.Error

func NewError(message string) *Error {
	return errors.New("PostingError", message)
}

func NewNotFoundError() *Error {
	return NewError("Такого поста не существует")
}
