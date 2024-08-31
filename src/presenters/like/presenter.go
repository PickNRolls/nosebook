package presenterlike

import (
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	presenterdto "nosebook/src/presenters/dto"
	"nosebook/src/services/auth"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	db            *sqlx.DB
	qb            squirrel.StatementBuilderType
	userPresenter UserPresenter
	resource      Resource

	err               *errors.Error
	auth              *auth.Auth
	linkedResourceIds uuid.UUIDs
	userIdsToFetchMap map[uuid.UUID]struct{}
	userIdsToFetch    uuid.UUIDs
	resourceLikersMap map[uuid.UUID]uuid.UUIDs
	users             []*presenterdto.User
	usersMap          usersMap
	out               likesMap
}

func New(db *sqlx.DB, userPresenter UserPresenter, resource Resource) *Presenter {
	out := &Presenter{
		db:            db,
		qb:            postgres.NewSquirrel(),
		userPresenter: userPresenter,
		resource:      resource,
	}

	out.reset()

	return out
}

func (this *Presenter) FindByResourceIds(ids uuid.UUIDs, auth *auth.Auth) (likesMap, *errors.Error) {
	this.linkedResourceIds = ids
	this.auth = auth

	this.fetchLikes()
	this.fetchLikeAdditionals()
	this.fetchUsers()

	return this.output()
}

func (this *Presenter) fetchLikes() {
	if this.err != nil {
		return
	}

	idColumn := this.resource.IDColumn()
	whereEq := squirrel.Eq{}
	whereEq[idColumn] = this.linkedResourceIds

	sql, args, _ := this.qb.
		Select("sub."+idColumn+" as resource_id", "sub.user_id").
		FromSelect(
			this.qb.
				Select(
					idColumn,
					"user_id",
					"row_number() over(partition by "+idColumn+") as row_number",
				).
				From(this.resource.Table()).
				Where(whereEq),
			"sub",
		).
		Where("sub.row_number <= 5").
		ToSql()

	likeDests := []*dest{}
	err := this.db.Select(&likeDests, sql, args...)
	if err != nil {
		this.err = errorFrom(err)
		return
	}

	for _, like := range likeDests {
		if _, has := this.userIdsToFetchMap[like.UserId]; !has {
			this.userIdsToFetch = append(this.userIdsToFetch, like.UserId)
			this.userIdsToFetchMap[like.UserId] = struct{}{}
		}

		if _, has := this.resourceLikersMap[like.ResourceId]; !has {
			this.resourceLikersMap[like.ResourceId] = make(uuid.UUIDs, 0)
		}
		this.resourceLikersMap[like.ResourceId] = append(this.resourceLikersMap[like.ResourceId], like.UserId)
	}
}

func (this *Presenter) fetchLikeAdditionals() {
	if this.err != nil {
		return
	}

	idColumn := this.resource.IDColumn()
	table := this.resource.Table()
	whereEq := squirrel.Eq{}
	whereEq[idColumn] = this.linkedResourceIds

	sql, args, _ := this.qb.
		Select("a."+idColumn+" as resource_id", "count", "b.user_id IS NOT NULL AS liked").
		FromSelect(
			this.qb.
				Select(idColumn, "count(*)").
				From(table).
				Where(whereEq).
				GroupBy(idColumn),
			"a",
		).
		JoinClause("LEFT OUTER JOIN "+table+" AS b").
		Suffix("ON a."+idColumn+" = b."+idColumn+" AND b.user_id = ?", this.auth.UserId).
		ToSql()

	additional := []struct {
		ResourceId uuid.UUID `db:"resource_id"`
		Count      int       `db:"count"`
		Liked      bool      `db:"liked"`
	}{}
	err := this.db.Select(&additional, sql, args...)
	if err != nil {
		this.err = errorFrom(err)
		return
	}

	for _, resourceId := range this.linkedResourceIds {
		this.out[resourceId] = &presenterdto.Likes{
			RandomFiveLikers: make([]*presenterdto.User, 0),
		}
	}

	for _, add := range additional {
		this.out[add.ResourceId].Count = add.Count
		this.out[add.ResourceId].Liked = add.Liked
	}
}

func (this *Presenter) fetchUsers() {
	if this.err != nil {
		return
	}

	this.users, this.err = this.userPresenter.FindByIds(this.userIdsToFetch)
	if this.err != nil {
		return
	}

	this.usersMap = usersMap{}
	for _, user := range this.users {
		this.usersMap[user.Id] = user
	}
}

func (this *Presenter) output() (likesMap, *errors.Error) {
	if this.err != nil {
		err := this.err
		this.reset()
		return nil, err
	}

	for resourceId, likerIds := range this.resourceLikersMap {
		for _, userId := range likerIds {
			this.out[resourceId].RandomFiveLikers = append(this.out[resourceId].RandomFiveLikers, this.usersMap[userId])
		}
	}

	out := this.out
	this.reset()
	return out, nil
}

func (this *Presenter) reset() {
	this.err = nil
	this.auth = nil
	this.linkedResourceIds = make(uuid.UUIDs, 0)
	this.userIdsToFetch = make(uuid.UUIDs, 0)
	this.userIdsToFetchMap = make(map[uuid.UUID]struct{})
	this.resourceLikersMap = make(map[uuid.UUID]uuid.UUIDs)
	this.users = make([]*presenterdto.User, 0)
	this.usersMap = make(usersMap)
	this.out = make(map[uuid.UUID]*presenterdto.Likes)
}
