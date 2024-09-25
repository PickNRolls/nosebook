package presentermessage

import (
	"context"
	"nosebook/src/errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Presenter struct {
	db            *sqlx.DB
	userPresenter UserPresenter
  tracer trace.Tracer
  findByIdsQuery *findByIdsQuery
}

func New(db *sqlx.DB, userPresenter UserPresenter) *Presenter {
	return &Presenter{
		db:            db,
		userPresenter: userPresenter,
    tracer: noop.Tracer{},
    findByIdsQuery: newFindByIdsQuery(db, userPresenter, noop.Tracer{}),
	}
}

func (this *Presenter) WithTracer(tracer trace.Tracer) *Presenter {
  this.tracer = tracer
  this.findByIdsQuery.tracer = tracer

  return this
}

func (this *Presenter) FindByIds(parent context.Context, ids []uuid.UUID) (map[uuid.UUID]*message, *errors.Error) {
  return this.findByIdsQuery.Do(parent, ids)
}

func (this *Presenter) Run() {
  this.findByIdsQuery.Run()
}

func (this *Presenter) OnDone() {
  this.findByIdsQuery.OnDone()
}

