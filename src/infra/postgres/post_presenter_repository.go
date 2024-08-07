package postgres

import (
	"nosebook/src/presenters/post_presenter/dto"
	"nosebook/src/presenters/post_presenter/interfaces"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostPresenterRepository struct {
	db *sqlx.DB
}

func NewPostPresenterRepository(db *sqlx.DB) interfaces.PostRepository {
	return &PostPresenterRepository{
		db: db,
	}
}

func findUsers(repo *PostPresenterRepository, ids []uuid.UUID) ([]*dto.UserDTO, error) {
	if len(ids) == 0 {
		return make([]*dto.UserDTO, 0), nil
	}

	var users []*dto.UserDTO
	query, args, err := sqlx.In(`SELECT
		id,
		first_name,
		last_name,
		nick,
		created_at
		  FROM users WHERE
		id IN (?)
	`, ids)
	if err != nil {
		return nil, err
	}

	query = repo.db.Rebind(query)

	err = repo.db.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *PostPresenterRepository) FindAuthors(ids []uuid.UUID) ([]*dto.UserDTO, error) {
	return findUsers(repo, ids)
}

func (repo *PostPresenterRepository) FindOwners(ids []uuid.UUID) ([]*dto.UserDTO, error) {
	return findUsers(repo, ids)
}
func (repo *PostPresenterRepository) FindLikers(ids []uuid.UUID) ([]*dto.UserDTO, error) {
	return findUsers(repo, ids)
}
