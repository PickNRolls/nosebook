package postgres

import (
	"nosebook/src/domain/sessions"
	"nosebook/src/services/user_authentication/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) interfaces.SessionRepository {
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
		user_id = $1
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
		session_id = $1
	`, id)

	if err != nil {
		return nil
	}

	return &session
}
