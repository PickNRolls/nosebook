package presenterpost

import (
	"nosebook/src/services/auth"

	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	userPresenter     userPresenter
	commentPresenter  commentPresenter
	likePresenter     likePresenter
	findByFilterQuery *findByFilterQuery
}

func New(db *sqlx.DB, userPresenter userPresenter, commentPresenter commentPresenter, likePresenter likePresenter) *Presenter {
	return &Presenter{
		userPresenter:     userPresenter,
		commentPresenter:  commentPresenter,
		likePresenter:     likePresenter,
		findByFilterQuery: newFindByFilterQuery(db),
	}
}

func (this *Presenter) FindByFilter(input *FindByFilterInput, a *auth.Auth) *FindByFilterOutput {
	return this.findByFilterQuery.FindByFilter(input, a, this.userPresenter, this.commentPresenter, this.likePresenter)
}
