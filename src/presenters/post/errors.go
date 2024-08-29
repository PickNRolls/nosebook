package presenterpost

import "nosebook/src/errors"

func newError(message string) *errors.Error {
	return errors.New("Post Presenter Error", message)
}
