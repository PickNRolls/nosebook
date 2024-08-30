package presentercomment

import "nosebook/src/errors"

func newError(message string) *errors.Error {
	return errors.New("Comment Presenter Error", message)
}

func errorFrom(err error) *errors.Error {
	return newError(err.Error())
}
