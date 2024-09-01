package permissionspost

import "nosebook/src/errors"

type Error = errors.Error

func newError(message string) *Error {
	return errors.New("Post Permissions Error", message)
}
