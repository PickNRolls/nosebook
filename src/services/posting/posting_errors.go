package posting

import "nosebook/src/errors"

type PostingError = errors.Error

func NewError(message string) *PostingError {
	return errors.New("PostingError", message)
}

func NewNotFoundError() *PostingError {
	return NewError("Такого поста не существует")
}
