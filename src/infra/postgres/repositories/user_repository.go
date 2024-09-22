package repos

import (
	"nosebook/src/domain/user"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db  *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db:  db,
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

	return &dest
}
