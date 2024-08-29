package domainpost

import "nosebook/src/errors"

type PostError = errors.Error

func NewError(message string) *PostError {
	return errors.New("PostError", message)
}
