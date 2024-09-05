package rootconvservice

import (
	querybuilder "nosebook/src/infra/query_builder"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func newUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (this *userRepository) Exists(id uuid.UUID) bool {
	qb := querybuilder.New()
	sql, args, _ := qb.Select("id").
		From("users").
		Where("id = ?", id).
		ToSql()

	dest := struct {
		Id uuid.UUID `db:"id"`
	}{}
	err := this.db.Get(&dest, sql, args...)
	if err != nil {
		return false
	}

	return true
}
