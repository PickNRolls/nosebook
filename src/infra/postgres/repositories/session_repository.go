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
	cache          Cache
	bufferUpdate   *worker.Buffer[*sessions.Session, error]
	bufferFindById *worker.Buffer[uuid.UUID, map[uuid.UUID]*sessions.Session]
}

type Cache interface {
	Set(id uuid.UUID, session *sessions.Session)
	Get(id uuid.UUID) (*sessions.Session, bool)
	Remove(id uuid.UUID)
}

type noopCache struct{}

func (this *noopCache) Set(id uuid.UUID, session *sessions.Session) {}
func (this *noopCache) Get(id uuid.UUID) (*sessions.Session, bool)  { return nil, false }
func (this *noopCache) Remove(id uuid.UUID)                         {}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	out := &SessionRepository{
		db:    db,
		cache: &noopCache{},
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

func (this *SessionRepository) CacheWith(cache Cache) *SessionRepository {
	this.cache = cache
	return this
}

func (this *SessionRepository) Create(session *sessions.Session) (*sessions.Session, error) {
	_, err := this.db.NamedExec(`INSERT INTO user_sessions (
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

	this.cache.Set(session.SessionId, session)

	return session, nil
}

func (this *SessionRepository) Remove(id uuid.UUID) (*sessions.Session, error) {
	var session sessions.Session
	err := this.db.Get(&session, `DELETE FROM user_sessions WHERE
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

	this.cache.Remove(session.SessionId)

	return &session, nil
}

func (this *SessionRepository) Update(session *sessions.Session) (*sessions.Session, error) {
	if _, has := this.cache.Get(session.SessionId); has {
		this.cache.Set(session.SessionId, session)
		go func() {
			this.bufferUpdate.Send(session)
		}()
		return session, nil
	}

	err := this.bufferUpdate.Send(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (this *SessionRepository) FindById(id uuid.UUID) *sessions.Session {
	if session, has := this.cache.Get(id); has {
		return session
	}

	m := this.bufferFindById.Send(id)
	session := m[id]
	this.cache.Set(id, session)
	return session
}
