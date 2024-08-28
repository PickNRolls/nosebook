package errors

import (
	"nosebook/src/errors"
)

type AuthenticatedError = errors.Error

func NewAuthenticatedError() *AuthenticatedError {
	return errors.New("Authenticated", "Only unauthenticated users can do it")
}
