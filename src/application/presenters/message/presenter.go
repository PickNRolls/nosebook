package presentermessage

import (
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
