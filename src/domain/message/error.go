package domainmessage

import "nosebook/src/errors"

type Error = errors.Error

func newError(message string) *Error {
	return errors.New("Message Error", message)
}
