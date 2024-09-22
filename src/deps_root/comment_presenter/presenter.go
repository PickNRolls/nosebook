package rootcommentpresenter

import (
	presentercomment "nosebook/src/application/presenters/comment"
	presentercommentlike "nosebook/src/application/presenters/comment_like"
	presenteruser "nosebook/src/application/presenters/user"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

func New(db *sqlx.DB, tracer trace.Tracer) *presentercomment.Presenter {
	userPresenter := presenteruser.New(db).WithTracer(tracer)
	likePresenter := presentercommentlike.New(db, userPresenter, tracer)

	presenter := presentercomment.New(
		db,
		likePresenter,
		userPresenter,
		newPermissions(db),
	).WithTracer(tracer)

	return presenter
}
