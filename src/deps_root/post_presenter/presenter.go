package rootpostpresenter

import (
	presenterpost "nosebook/src/application/presenters/post"
	"nosebook/src/application/presenters/post_like"
	presenteruser "nosebook/src/application/presenters/user"
	rootcommentpresenter "nosebook/src/deps_root/comment_presenter"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

func New(db *sqlx.DB, tracer trace.Tracer) *presenterpost.Presenter {
	userPresenter := presenteruser.New(db).WithTracer(tracer)
	likePresenter := presenterpostlike.New(db, userPresenter).WithTracer(tracer)
	presenter := presenterpost.New(
		db,
		userPresenter,
		rootcommentpresenter.New(db).WithTracer(tracer),
		likePresenter,
		&permissions{},
	).WithTracer(tracer)

	return presenter
}
