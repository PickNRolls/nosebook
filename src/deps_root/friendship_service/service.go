package rootfriendshipservice

import (
	"nosebook/src/application/services/friendship"

	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *friendship.Service {
	service := friendship.New(newRepository(db))

	return service
}
