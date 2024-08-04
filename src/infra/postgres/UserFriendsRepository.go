package postgres

import (
	"nosebook/src/domain/friendship"
	"nosebook/src/services/friendship/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserFriendsRepository struct {
	db *sqlx.DB
}

func NewUserFriendsRepository(db *sqlx.DB) interfaces.UserFriendsRepository {
	return &UserFriendsRepository{
		db: db,
	}
}

func (repo *UserFriendsRepository) Create(request *friendship.FriendRequest) (*friendship.FriendRequest, error) {
	_, err := repo.db.NamedExec(`INSERT INTO friendship_requests (
	  requester_id,
	  responder_id,
	  message,
	  accepted,
	  viewed,
	  created_at
	) VALUES (
	  :requester_id,
	  :responder_id,
	  :message,
	  :accepted,
	  :viewed,
	  :created_at
	)`, request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (repo *UserFriendsRepository) Update(request *friendship.FriendRequest) (*friendship.FriendRequest, error) {
	_, err := repo.db.NamedExec(`UPDATE friendship_requests SET
		accepted = :accepted,
		viewed = :viewed
			WHERE
		requester_id = :requester_id AND
		responder_id = :responder_id
	`, request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func (repo *UserFriendsRepository) FindByBoth(requesterId uuid.UUID, responderId uuid.UUID) *friendship.FriendRequest {
	request := friendship.FriendRequest{}
	err := repo.db.Get(&request, `SELECT
		requester_id,
		responder_id,
		message,
		accepted,
		viewed,
		created_at
			FROM friendship_requests WHERE
		requester_id = $1 AND responder_id = $2
	`, requesterId, responderId)

	if err != nil {
		return nil
	}

	return &request
}
