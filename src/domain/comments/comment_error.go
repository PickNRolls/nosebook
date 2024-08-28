package comments

import "nosebook/src/errors"

type CommentError = errors.Error

func NewError(message string) *CommentError {
	return errors.New("CommentError", message)
}
