package presenterpost

import (
	"nosebook/src/application/services/auth"

	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	userPresenter     UserPresenter
	commentPresenter  CommentPresenter
	likePresenter     LikePresenter
	findByFilterQuery *findByFilterQuery
}

func New(
	db *sqlx.DB,
	userPresenter UserPresenter,
	commentPresenter CommentPresenter,
	likePresenter LikePresenter,
	permissions Permissions,
) *Presenter {
	return &Presenter{
		userPresenter:    userPresenter,
		commentPresenter: commentPresenter,
		likePresenter:    likePresenter,
		findByFilterQuery: newFindByFilterQuery(
			db,
			userPresenter,
			commentPresenter,
			likePresenter,
			permissions,
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
