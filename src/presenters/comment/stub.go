//go:build exclude

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

func (this *Presenter) FindByPostId(id uuid.UUID) *FindByFilterOutput {
	return nil
}
