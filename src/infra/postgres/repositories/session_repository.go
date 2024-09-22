package repos

import (
	userauth "nosebook/src/application/services/user_auth"
	"nosebook/src/domain/sessions"
	"nosebook/src/lib/cache"
	"nosebook/src/lib/clock"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db  *sqlx.DB
	lru *cache.LRU[string, *sessions.Session]
}

func NewSessionRepository(db *sqlx.DB) userauth.SessionRepository {
	return &SessionRepository{
		db:  db,
		lru: cache.NewLRU[string, *sessions.Session](2048),
	}
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
	_, has := this.lru.Get(session.SessionId.String())
	if has {
		this.lru.Set(session.SessionId.String(), session)
		return session, nil
	}

	_, err := this.db.NamedExec(`UPDATE user_sessions SET
		expires_at = :expires_at
			WHERE
		user_id = :user_id AND session_id = :session_id
	`, session)
	if err != nil {
		return nil, err
	}

  this.lru.Set(session.SessionId.String(), session)

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
	cached, has := this.lru.Get(id.String())
	if has {
		return cached
	}

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

	this.lru.Set(id.String(), &dest)

	return &dest
}
