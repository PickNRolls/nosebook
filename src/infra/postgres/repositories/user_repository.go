package repos

import (
	prometheusmetrics "nosebook/src/deps_root/worker"
	"nosebook/src/domain/user"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/worker"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db             *sqlx.DB
	bufferUpdate   *worker.Buffer[*bufferUpdate, error]
	bufferFindById *worker.Buffer[uuid.UUID, map[uuid.UUID]*domainuser.User]
}

type bufferUpdate struct {
	userId    uuid.UUID
	timestamp time.Time
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	out := &UserRepository{
		db: db,
	}

	out.bufferUpdate = worker.NewBuffer(func(updates []*bufferUpdate) error {
		sql := `UPDATE users as u SET last_activity_at = v.last_activity_at
    FROM (VALUES `
		args := []any{}

		for i, update := range updates {
			last := i == len(updates)-1
			argNum := len(args) + 1
			suffix := "($" + strconv.Itoa(argNum) + "::uuid, $" + strconv.Itoa(argNum+1) + "::timestamp)"
			if !last {
				suffix += ","
			}

			sql += suffix
			args = append(args, update.userId, update.timestamp)
		}

		sql += ") v(id, last_activity_at) WHERE u.id = v.id"
		_, err := db.Exec(sql, args...)
		return err
	}, prometheusmetrics.UsePrometheusMetrics("users_update"))

	qb := querybuilder.New()
	out.bufferFindById = worker.NewBuffer(func(ids []uuid.UUID) map[uuid.UUID]*domainuser.User {
		out := map[uuid.UUID]*domainuser.User{}

		dests := []*domainuser.User{}
		sql, args, _ := qb.Select("id", "first_name", "last_name", "nick", "passhash", "created_at", "last_activity_at").
			From("users").
			Where(squirrel.Eq{"id": ids}).
			ToSql()

		err := db.Select(&dests, sql, args...)
		if err != nil {
			return out
		}

		for _, dest := range dests {
			out[dest.Id] = dest
		}

		return out
	}, prometheusmetrics.UsePrometheusMetrics("users_find"))

	return out
}

func (this *UserRepository) Run() {
	go this.bufferUpdate.Run()
	go this.bufferFindById.Run()
}

func (this *UserRepository) OnDone() {
	this.bufferUpdate.Stop()
	this.bufferFindById.Stop()
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
	return this.bufferUpdate.Send(&bufferUpdate{
		userId:    userId,
		timestamp: t,
	})
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
	m := this.bufferFindById.Send(id)
	return m[id]
}
