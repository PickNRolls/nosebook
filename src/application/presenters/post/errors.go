package presenterpost

import "nosebook/src/errors"

func newError(message string) *errors.Error {
	return errors.New("Post Presenter Error", message)
}

func errorFrom(err error) *errors.Error {
	return newError(err.Error())
}
