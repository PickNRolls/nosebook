package errors

import (
	"nosebook/src/errors"
)

type AuthorizedError = errors.Error

func NewAuthorizedError() *AuthorizedError {
	return errors.New("Authorized", "Only unauthorized users can do it")
}
