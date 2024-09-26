package repos

import (
	prometheusmetrics "nosebook/src/deps_root/worker"
	"nosebook/src/domain/sessions"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/clock"
	"nosebook/src/lib/worker"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db             *sqlx.DB
	bufferUpdate   *worker.Buffer[*sessions.Session, error]
	bufferFindById *worker.Buffer[uuid.UUID, map[uuid.UUID]*sessions.Session]
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	out := &SessionRepository{
		db: db,
	}

	out.bufferUpdate = worker.NewBuffer(func(sessions []*sessions.Session) error {
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
		return err
	}, prometheusmetrics.UsePrometheusMetrics("sessions_update"))

	out.bufferFindById = worker.NewBuffer(func(ids []uuid.UUID) map[uuid.UUID]*sessions.Session {
		out := map[uuid.UUID]*sessions.Session{}
		qb := querybuilder.New()

		dests := []*sessions.Session{}
		sql, args, _ := qb.Select("session_id", "user_id", "created_at", "expires_at").
			From("user_sessions").
			Where(squirrel.Eq{"session_id": ids}).
			Where("expires_at > ?", clock.Now()).
			ToSql()
		err := db.Select(&dests, sql, args...)
		if err != nil {
			return out
		}

		for _, dest := range dests {
			out[dest.SessionId] = dest
		}

		return out
	}, prometheusmetrics.UsePrometheusMetrics("sessions_find"))

	return out
}

func (this *SessionRepository) Run() {
	go this.bufferUpdate.Run()
	go this.bufferFindById.Run()
}

func (this *SessionRepository) OnDone() {
	this.bufferUpdate.Stop()
	this.bufferFindById.Stop()
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
	err := this.bufferUpdate.Send(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (this *SessionRepository) FindById(id uuid.UUID) *sessions.Session {
	m := this.bufferFindById.Send(id)
	return m[id]
}
