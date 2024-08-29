package rootfriendshipservice

import (
	"nosebook/src/services/friendship"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *friendship.Service {
	service := friendship.New(newRepository(db))

	return service
}
