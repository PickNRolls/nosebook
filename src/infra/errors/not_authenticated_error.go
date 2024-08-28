package errors

import (
	"nosebook/src/errors"
)

type NotAuthenticatedError = errors.Error

func NewNotAuthenticatedError() *NotAuthenticatedError {
	return errors.New("Not Authenticated", "You are not authenticated")
}
