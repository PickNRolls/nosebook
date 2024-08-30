package rootcommentpresenter

import (
	presentercomment "nosebook/src/presenters/comment"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *presentercomment.Presenter {
	presenter := presentercomment.New(db)

	return presenter
}
