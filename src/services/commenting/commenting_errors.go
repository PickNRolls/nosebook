package commenting

import "nosebook/src/errors"

type CommentError = errors.Error

func NewError(message string) *CommentError {
	return errors.New("CommentError", message)
}

func NewPostNotFoundError() *CommentError {
	return NewError("Такого поста не существует")
}
