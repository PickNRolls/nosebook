package errors

import (
	"nosebook/src/errors"
)

type NotAuthorizedError = errors.Error

func NewNotAuthorizedError() *NotAuthorizedError {
	return errors.New("NotAuthorized", "You are not authorized")
}
