package conversation

import "nosebook/src/errors"

type Error = errors.Error

func newError(message string) *Error {
	return errors.New("Conversation Error", message)
}
