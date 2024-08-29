package repos

import (
	"nosebook/src/domain/sessions"
	userauth "nosebook/src/services/user_auth"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) userauth.SessionRepository {
	return &SessionRepository{
		db: db,
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

func (repo *SessionRepository) Update(session *sessions.Session) (*sessions.Session, error) {
	_, err := repo.db.NamedExec(`UPDATE user_sessions SET
		expires_at = :expires_at
			WHERE
		user_id = :user_id AND session_id = :session_id
	`, session)
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
		user_id = $1 AND expires_at > NOW()
	`, userId)

	if err != nil {
		return nil
	}

	return &session
}

func (repo *SessionRepository) FindById(id uuid.UUID) *sessions.Session {
	session := sessions.Session{}
	err := repo.db.Get(&session, `SELECT
		session_id,
		user_id,
		created_at,
		expires_at
			FROM user_sessions WHERE
		session_id = $1 AND expires_at > NOW()
	`, id)

	if err != nil {
		return nil
	}

	return &session
}
