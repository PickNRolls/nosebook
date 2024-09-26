package rootconvservice

import (
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/worker"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db     *sqlx.DB
	buffer *worker.Buffer[uuid.UUID, map[uuid.UUID]struct{}]
}

func newUserRepository(db *sqlx.DB) *userRepository {
	qb := querybuilder.New()

	return &userRepository{
		db: db,
		buffer: worker.NewBuffer(func(ids []uuid.UUID) map[uuid.UUID]struct{} {
			out := map[uuid.UUID]struct{}{}

			sql, args, _ := qb.Select("id").
				From("users").
				Where(squirrel.Eq{"id": ids}).
				ToSql()

			dests := []struct {
				Id uuid.UUID `db:"id"`
			}{}
			err := db.Select(&dests, sql, args...)
			if err != nil {
				return out
			}

			for _, dest := range dests {
				out[dest.Id] = struct{}{}
			}

			return out
		}),
	}
}

func (this *userRepository) Run() {
	this.buffer.Run()
}

func (this *userRepository) OnDone() {
	this.buffer.Stop()
}

func (this *userRepository) Exists(id uuid.UUID) bool {
	m := this.buffer.Send(id)
	_, has := m[id]

	return has
}
