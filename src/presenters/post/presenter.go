package presenterpost

import (
	"nosebook/src/services/auth"

	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	findByFilterQuery *findByFilterQuery
}

func New(db *sqlx.DB) *Presenter {
	return &Presenter{
		findByFilterQuery: newFindByFilterQuery(db),
	}
}

func (this *Presenter) FindByFilter(input *FindByFilterInput, a *auth.Auth) *FindByFilterOutput {
	return this.findByFilterQuery.FindByFilter(input, a)
}

