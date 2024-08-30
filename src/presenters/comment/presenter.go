package presentercomment

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Presenter {
	return &Presenter{
		db: db,
	}
}

func (this *Presenter) FindByFilter(input *FindByFilterInput) *FindByFilterOutput {
	return nil
}

func (this *Presenter) FindByPostId(id uuid.UUID) *FindByFilterOutput {
	return this.FindByFilter(&FindByFilterInput{
		PostId: id.String(),
	})
}
