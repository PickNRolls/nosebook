package repos

import (
	prometheusmetrics "nosebook/src/deps_root/worker"
	"nosebook/src/domain/user"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/worker"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db             *sqlx.DB
	cache          Cache
	bufferUpdate   *worker.Buffer[*domainuser.User, error]
	bufferFindById *worker.Buffer[uuid.UUID, map[uuid.UUID]*domainuser.User]
}

type Cache interface {
	Set(id uuid.UUID, session *domainuser.User)
	Get(id uuid.UUID) (*domainuser.User, bool)
	Remove(id uuid.UUID)
}

type noopCache struct{}

func (this *noopCache) Set(id uuid.UUID, session *domainuser.User) {}
func (this *noopCache) Get(id uuid.UUID) (*domainuser.User, bool)  { return nil, false }
func (this *noopCache) Remove(id uuid.UUID)                        {}

func New(db *sqlx.DB) *UserRepository {
	out := &UserRepository{
		db:    db,
		cache: &noopCache{},
	}

	out.bufferUpdate = worker.NewBuffer(func(users []*domainuser.User) error {
		sql := `UPDATE users as u
    SET
      last_activity_at = v.last_activity_at,
      avatar_url = v.avatar_url,
      avatar_updated_at = v.avatar_updated_at
    FROM (VALUES `
		args := []any{}

		for i, user := range users {
			last := i == len(users)-1
			argNum := len(args) + 1
			suffix := "($" + strconv.Itoa(argNum) + "::uuid, $" + strconv.Itoa(argNum+1) + "::timestamp, $" + strconv.Itoa(argNum+2) + ", $" + strconv.Itoa(argNum+3) + "::timestamp)"
			if !last {
				suffix += ","
			}

			sql += suffix
			args = append(args, user.Id, user.LastActivityAt, user.AvatarUrl, user.AvatarUpdatedAt)
		}

		sql += ") v(id, last_activity_at, avatar_url, avatar_updated_at) WHERE u.id = v.id"
		_, err := db.Exec(sql, args...)
		return err
	}, prometheusmetrics.UsePrometheusMetrics("users_update"))

	qb := querybuilder.New()
	out.bufferFindById = worker.NewBuffer(func(ids []uuid.UUID) map[uuid.UUID]*domainuser.User {
		out := map[uuid.UUID]*domainuser.User{}

		dests := []*domainuser.User{}
		sql, args, _ := qb.Select("id", "first_name", "last_name", "nick", "passhash", "created_at", "avatar_url", "avatar_updated_at", "last_activity_at").
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

func (this *UserRepository) CacheWith(cache Cache) *UserRepository {
	this.cache = cache
	return this
}

func (repo *UserRepository) Create(user *domainuser.User) (*domainuser.User, error) {
	_, err := repo.db.NamedExec(`INSERT INTO users (
		id,
	  first_name,
	  last_name,
	  passhash,
	  nick,
	  created_at,
    avatar_url,
	  last_activity_at
	) VALUES (
		:id,
	  :first_name,
	  :last_name,
	  :passhash,
	  :nick,
	  :created_at,
    :avatar_url,
	  :last_activity_at
	) RETURNING (
		id,
		first_name,
		last_name,
		passhash,
		nick,
		created_at,
    avatar_url,
		last_activity_at
	)`, user)
	if err != nil {
		return nil, err
	}

	return user, nil
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
		last_activity_at,
    avatar_url,
    avatar_updated_at
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
    avatar_url,
    avatar_updated_at,
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
	if user, has := this.cache.Get(id); has {
		return user
	}

	m := this.bufferFindById.Send(id)
	user := m[id]

	this.cache.Set(id, user)

	return user
}

func (this *UserRepository) Save(user *domainuser.User) *errors.Error {
	if _, has := this.cache.Get(user.Id); has {
		this.cache.Set(user.Id, user)
		go func() {
			this.bufferUpdate.Send(user)
		}()
		return nil
	}

	return errors.From(this.bufferUpdate.Send(user))
}
