package presentermessage

import (
	"context"
	prometheusmetrics "nosebook/src/deps_root/worker"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/worker"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

type findByIdsQuery struct {
	db            *sqlx.DB
	userPresenter UserPresenter
	buffer        *worker.Buffer[[]uuid.UUID, *findByIdsQueryBufferOut]
	done          chan struct{}
	tracer        trace.Tracer
}

type findByIdsQueryBufferOut struct {
	err  *errors.Error
	data map[uuid.UUID]*message
}

func newFindByIdsQuery(db *sqlx.DB, userPresenter UserPresenter, tracer trace.Tracer) *findByIdsQuery {
	return &findByIdsQuery{
		db:            db,
		userPresenter: userPresenter,
		done:          make(chan struct{}),
		tracer:        tracer,
	}
}

func (this *findByIdsQuery) Do(parent context.Context, ids uuid.UUIDs) (map[uuid.UUID]*message, *errors.Error) {
	_, span := this.tracer.Start(parent, "message_presenter.find_by_ids")
	defer span.End()

	out := this.buffer.Send(ids)
	filtered := map[uuid.UUID]*message{}
	for _, id := range ids {
		filtered[id] = out.data[id]
	}

	return filtered, out.err
}

func (this *findByIdsQuery) Run() {
	this.buffer = worker.NewBuffer(func(values [][]uuid.UUID) *findByIdsQueryBufferOut {
		unique := map[uuid.UUID]struct{}{}
		for _, row := range values {
			for _, id := range row {
				if _, has := unique[id]; !has {
					unique[id] = struct{}{}
				}
			}
		}

		ids := []uuid.UUID{}
		for id := range unique {
			ids = append(ids, id)
		}

		qb := querybuilder.New()
		sql, args, _ := qb.
			Select("id", "author_id", "text", "chat_id", "reply_to", "created_at").
			From("messages").
			Where(squirrel.Eq{"id": ids}).
			ToSql()

		dests := []*dest{}
		err := errors.From(this.db.Select(&dests, sql, args...))
		if err != nil {
			return &findByIdsQueryBufferOut{
				err:  err,
				data: nil,
			}
		}

		data, err := mapDests(context.TODO(), this.userPresenter, dests)
		return &findByIdsQueryBufferOut{
			err:  err,
			data: data,
		}
	}, prometheusmetrics.UsePrometheusMetrics("message_find"))

	this.buffer.Run()
}

func (this *findByIdsQuery) OnDone() {
	this.done <- struct{}{}
	close(this.done)
}
