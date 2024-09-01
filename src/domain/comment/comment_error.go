package domaincomment

import "nosebook/src/errors"

type Error = errors.Error

func NewError(message string) *Error {
	return errors.New("Comment Error", message)
}
