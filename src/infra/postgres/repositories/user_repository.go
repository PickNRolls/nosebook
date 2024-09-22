package repos

import (
	userauth "nosebook/src/application/services/user_auth"
	"nosebook/src/domain/user"
	"nosebook/src/lib/cache"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db  *sqlx.DB
	lru *cache.LRU[string, *domainuser.User]
}

func NewUserRepository(db *sqlx.DB) userauth.UserRepository {
	return &UserRepository{
		db:  db,
		lru: cache.NewLRU[string, *domainuser.User](2048),
	}
}

func (repo *UserRepository) Create(user *domainuser.User) (*domainuser.User, error) {
	_, err := repo.db.NamedExec(`INSERT INTO users (
		id,
	  first_name,
	  last_name,
	  passhash,
	  nick,
	  created_at,
	  last_activity_at
	) VALUES (
		:id,
	  :first_name,
	  :last_name,
	  :passhash,
	  :nick,
	  :created_at,
	  :last_activity_at
	) RETURNING (
		id,
		first_name,
		last_name,
		passhash,
		nick,
		created_at,
		last_activity_at
	)`, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (this *UserRepository) UpdateActivity(userId uuid.UUID, t time.Time) error {
	if cached, has := this.lru.Get(userId.String()); has {
		cached.LastActivityAt = t
		return nil
	}

	_, err := this.db.Exec(`UPDATE users SET
		last_activity_at = $1
			WHERE
		id = $2
	`, t, userId)

	return err
}

func (this *UserRepository) FindAll() ([]*domainuser.User, error) {
	var all []*domainuser.User
	err := this.db.Select(&all, `SELECT
		id,
		first_name,
		last_name,
		nick,
		passhash,
		created_at,
		last_activity_at
			FROM users
	`)

	if err != nil {
		return make([]*domainuser.User, 0), nil
	}

	return all, nil
}

func (repo *UserRepository) FindByNick(nick string) *domainuser.User {
	user := domainuser.User{}
	err := repo.db.Get(&user, `SELECT
		id,
	  first_name,
	  last_name,
	  nick,
	  passhash,
	  created_at,
	  last_activity_at
	    FROM users WHERE
    nick = $1
	`, nick)

	if err != nil {
		return nil
	}

	return &user
}

func (this *UserRepository) FindById(id uuid.UUID) *domainuser.User {
	if cached, has := this.lru.Get(id.String()); has {
		return cached
	}

	dest := domainuser.User{}
	err := this.db.Get(&dest, `SELECT
		id,
		first_name,
		last_name,
		nick,
		passhash,
		created_at,
		last_activity_at
			FROM users WHERE
		id = $1
	`, id)

	if err != nil {
		return nil
	}

	this.lru.Set(id.String(), &dest)

	return &dest
}
