package presenterfriendship

import (
	"nosebook/src/errors"

	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	db            *sqlx.DB
	userPresenter UserPresenter
}

func New(db *sqlx.DB, userPresenter UserPresenter) *Presenter {
	return &Presenter{
		db:            db,
		userPresenter: userPresenter,
	}
}

func newError(msg string) *errors.Error {
	return errors.New("Friendship Presenter Error", msg)
}

func errorFrom(err error) *errors.Error {
	return newError(err.Error())
}
