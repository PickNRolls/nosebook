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
		userPresenter:    userPresenter,
		commentPresenter: commentPresenter,
		likePresenter:    likePresenter,
		findByFilterQuery: newFindByFilterQuery(
			db,
			userPresenter,
			commentPresenter,
			likePresenter,
		),
	}
}

func (this *Presenter) FindByFilter(input *FindByFilterInput, a *auth.Auth) *FindByFilterOutput {
	return this.findByFilterQuery.FindByFilter(input, a)
}

func (this *Presenter) FindById(id string, a *auth.Auth) *post {
	out := this.FindByFilter(&FindByFilterInput{
		Ids: []string{id},
	}, a)

	if out.Data != nil && len(out.Data) != 0 {
		return out.Data[0]
	}

	return nil
}
