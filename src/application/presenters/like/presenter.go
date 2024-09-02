package presenterlike

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	db            *sqlx.DB
	qb            squirrel.StatementBuilderType
	userPresenter UserPresenter
	resource      Resource
}

type additionalDest struct {
	ResourceId uuid.UUID `db:"resource_id"`
	Count      int       `db:"count"`
	Liked      bool      `db:"liked"`
}

func New(db *sqlx.DB, userPresenter UserPresenter, resource Resource) *Presenter {
	return &Presenter{
		db:            db,
		qb:            postgres.NewSquirrel(),
		userPresenter: userPresenter,
		resource:      resource,
	}
}

func (this *Presenter) FindByResourceIds(ids uuid.UUIDs, auth *auth.Auth) (likesMap, *errors.Error) {
	userIdsMap := map[uuid.UUID]struct{}{}
	userIds := []uuid.UUID{}
	userMap := map[uuid.UUID]*presenterdto.User{}
	resourceLikerIdsMap := map[uuid.UUID]uuid.UUIDs{}
	dests := []*dest{}

	err := func() error {
		idColumn := this.resource.IDColumn()
		whereEq := squirrel.Eq{}
		whereEq[idColumn] = ids

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

		err := this.db.Select(&dests, sql, args...)
		if err != nil {
			return err
		}

		for _, like := range dests {
			if _, has := userIdsMap[like.UserId]; !has {
				userIds = append(userIds, like.UserId)
				userIdsMap[like.UserId] = struct{}{}
			}

			if _, has := resourceLikerIdsMap[like.ResourceId]; !has {
				resourceLikerIdsMap[like.ResourceId] = make(uuid.UUIDs, 0)
			}
			resourceLikerIdsMap[like.ResourceId] = append(resourceLikerIdsMap[like.ResourceId], like.UserId)
		}

		users, error := this.userPresenter.FindByIds(userIds)
		if error != nil {
			return error.Unwrap()
		}

		for _, user := range users {
			userMap[user.Id] = user
		}

		return nil
	}()
	if err != nil {
		return nil, errorFrom(err)
	}

	additional, error := func() ([]*additionalDest, *errors.Error) {
		idColumn := this.resource.IDColumn()
		table := this.resource.Table()
		whereEq := squirrel.Eq{}
		whereEq[idColumn] = ids

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
			Suffix("ON a."+idColumn+" = b."+idColumn+" AND b.user_id = ?", auth.UserId).
			ToSql()

		additional := []*additionalDest{}
		err := this.db.Select(&additional, sql, args...)
		if err != nil {
			return nil, errorFrom(err)
		}

		return additional, nil
	}()
	if error != nil {
		return nil, error
	}

	out := func() likesMap {
		out := likesMap{}

		for _, resourceId := range ids {
			out[resourceId] = &presenterdto.Likes{
				RandomFiveLikers: make([]*presenterdto.User, 0),
			}
		}

		for _, add := range additional {
			out[add.ResourceId].Count = add.Count
			out[add.ResourceId].Liked = add.Liked
		}

		for resourceId, likerIds := range resourceLikerIdsMap {
			for _, userId := range likerIds {
				out[resourceId].RandomFiveLikers = append(out[resourceId].RandomFiveLikers, userMap[userId])
			}
		}

		return out
	}()

	return out, nil
}
