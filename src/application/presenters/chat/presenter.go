package presenterchat

import "github.com/jmoiron/sqlx"

type Presenter struct {
	db               *sqlx.DB
	userPresenter    UserPresenter
	messagePresenter MessagePresenter
}

func New(db *sqlx.DB, userPresenter UserPresenter, messagePresenter MessagePresenter) *Presenter {
	return &Presenter{
		db:               db,
		userPresenter:    userPresenter,
		messagePresenter: messagePresenter,
	}
}
