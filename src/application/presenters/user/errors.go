package presenteruser

import "nosebook/src/errors"

func newError(message string) *errors.Error {
	return errors.New("User Presenter Error", message)
}

func errorFrom(err error) *errors.Error {
	return newError(err.Error())
}
