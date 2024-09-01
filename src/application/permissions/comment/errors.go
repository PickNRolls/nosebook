package permissionscomment

import "nosebook/src/errors"

type Error = errors.Error

func newError(message string) *Error {
	return errors.New("Comment Permissions Error", message)
}
