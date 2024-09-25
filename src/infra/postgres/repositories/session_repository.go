package repos

import (
	"nosebook/src/domain/sessions"
	"nosebook/src/lib/clock"
	"nosebook/src/lib/worker"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

type SessionRepository struct {
	db     *sqlx.DB
	ticker *time.Ticker
	buffer *worker.Buffer[*sessions.Session, error, time.Time]
	done   chan struct{}
}

var SessionsInWorkerTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "app_sessions_in_worker_total",
		Help: "Total count of sessions waiting for update in worker queue",
	},
)
var SessionsInWorkerCurrent = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "app_sessions_in_worker_current",
		Help: "Current count of sessions waiting for update in worker queue",
	},
)
var SessionsInWorkerBatchSize = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "app_sessions_in_worker_batch_size",
		Help:    "Number of sessions in one batch unit of work",
		Buckets: prometheus.ExponentialBuckets(1, 2, 10),
	},
)
var SessionsInWorkerUnitElapsed = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "app_sessions_in_worker_unit_elapsed_seconds",
		Help:    "Elapsed seconds of unit of work of sessions update worker",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 12),
	},
)

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	out := &SessionRepository{
		db:     db,
		ticker: time.NewTicker(time.Millisecond * 20),
	}

	buffer := worker.NewBuffer(func(sessions []*sessions.Session) error {
		SessionsInWorkerCurrent.Set(0)
		SessionsInWorkerBatchSize.Observe(float64(len(sessions)))
		before := clock.Now()

		if len(sessions) == 0 {
			return nil
		}

		sql := `UPDATE user_sessions as u SET expires_at = v.expires_at
    FROM (VALUES `
		args := []any{}

		for i, session := range sessions {
			last := i == len(sessions)-1
			argNum := len(args) + 1
			suffix := "($" + strconv.Itoa(argNum) + "::uuid, $" + strconv.Itoa(argNum+1) + "::timestamp)"
			if !last {
				suffix += ","
			}

			sql += suffix
			args = append(args, session.SessionId, session.ExpiresAt)
		}

		sql += ") v(id, expires_at) WHERE u.session_id = v.id"

		_, err := db.Exec(sql, args...)
		after := clock.Now()
		SessionsInWorkerUnitElapsed.Observe(float64(after.Sub(before).Seconds()))
		return err
	}, out.ticker.C, out.done, 256, worker.FlushEmpty)

	out.buffer = buffer

	return out
}

func (this *SessionRepository) Run() {
	this.buffer.Run()
}

func (this *SessionRepository) OnDispose() {
	this.done <- struct{}{}
	close(this.done)

	this.ticker.Stop()
}

func (repo *SessionRepository) Create(session *sessions.Session) (*sessions.Session, error) {
	_, err := repo.db.NamedExec(`INSERT INTO user_sessions (
	  session_id,
	  user_id,
	  created_at,
	  expires_at
	) VALUES (
	  :session_id,
	  :user_id,
	  :created_at,
	  :expires_at
	)`, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (repo *SessionRepository) Remove(id uuid.UUID) (*sessions.Session, error) {
	var session sessions.Session
	err := repo.db.Get(&session, `DELETE FROM user_sessions WHERE
		session_id = $1
			RETURNING
		session_id,
		user_id,
		created_at,
		expires_at
	`, id)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (this *SessionRepository) Update(session *sessions.Session) (*sessions.Session, error) {
	SessionsInWorkerTotal.Inc()
	SessionsInWorkerCurrent.Inc()

	err := this.buffer.Send(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (repo *SessionRepository) FindByUserId(userId uuid.UUID) *sessions.Session {
	session := sessions.Session{}
	err := repo.db.Get(&session, `SELECT
		session_id,
		user_id,
		created_at,
		expires_at
			FROM user_sessions WHERE
		user_id = $1 AND expires_at > $2
	`, userId, clock.Now())

	if err != nil {
		return nil
	}

	return &session
}

func (this *SessionRepository) FindById(id uuid.UUID) *sessions.Session {
	dest := sessions.Session{}
	err := this.db.Get(&dest, `SELECT
		session_id,
		user_id,
		created_at,
		expires_at
			FROM user_sessions WHERE
		session_id = $1 AND expires_at > $2
	`, id, clock.Now())

	if err != nil {
		return nil
	}

	return &dest
}
