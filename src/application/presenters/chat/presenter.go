package presenterchat

import (
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Presenter struct {
	db               *sqlx.DB
	userPresenter    UserPresenter
	messagePresenter MessagePresenter
  tracer trace.Tracer
}

func New(db *sqlx.DB, userPresenter UserPresenter, messagePresenter MessagePresenter) *Presenter {
	return &Presenter{
		db:               db,
		userPresenter:    userPresenter,
		messagePresenter: messagePresenter,
    tracer: noop.Tracer{},
	}
}

func (this *Presenter) WithTracer(tracer trace.Tracer) *Presenter {
  this.tracer = tracer

  return this
}
