package postgres

import (
	"nosebook/src/domain/users"
	"nosebook/src/services/user_authentication/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) interfaces.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) Create(user *users.User) (*users.User, error) {
	_, err := repo.db.NamedExec(`INSERT INTO users (
		id,
	  first_name,
	  last_name,
	  passhash,
	  nick,
	  created_at
	) VALUES (
		:id,
	  :first_name,
	  :last_name,
	  :passhash,
	  :nick,
	  :created_at
	) RETURNING (
		id,
		first_name,
		last_name,
		passhash,
		nick,
		created_at
	)`, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo *UserRepository) FindByNick(nick string) *users.User {
	user := users.User{}
	err := repo.db.Get(&user, `SELECT
	  first_name,
	  last_name,
	  nick,
	  passhash,
	  created_at
	    FROM users WHERE
    nick = $1
	`, nick)

	if err != nil {
		return nil
	}

	return &user
}

func (repo *UserRepository) FindById(id uuid.UUID) *users.User {
	user := users.User{}
	err := repo.db.Get(&user, `SELECT
		id,
		first_name,
		last_name,
		nick,
		passhash,
		created_at
			FROM users WHERE
		id = $1
	`, id)

	if err != nil {
		return nil
	}

	return &user
}
