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
	userPresenter userPresenter

	err               *errors.Error
	auth              *auth.Auth
	postIds           uuid.UUIDs
	userIdsToFetchMap map[uuid.UUID]struct{}
	userIdsToFetch    uuid.UUIDs
	postLikersMap     map[uuid.UUID]uuid.UUIDs
	users             []*presenterdto.User
	usersMap          map[uuid.UUID]*presenterdto.User
	out               map[uuid.UUID]*presenterdto.Likes
}

type likeDest struct {
	PostId uuid.UUID `db:"post_id"`
	UserId uuid.UUID `db:"user_id"`
}

func New(db *sqlx.DB, userPresenter userPresenter) *Presenter {
	out := &Presenter{
		db:            db,
		qb:            postgres.NewSquirrel(),
		userPresenter: userPresenter,
	}

	out.reset()

	return out
}

func (this *Presenter) FindByPostIds(ids uuid.UUIDs, auth *auth.Auth) (map[uuid.UUID]*presenterdto.Likes, *errors.Error) {
	this.postIds = ids
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

	sql, args, _ := this.qb.
		Select("sub.post_id, sub.user_id").
		FromSelect(
			this.qb.
				Select(
					"post_id",
					"user_id",
					"row_number() over(partition by post_id) as row_number",
				).
				From("post_likes").
				Where(squirrel.Eq{"post_id": this.postIds}),
			"sub",
		).
		Where("sub.row_number <= 5").
		ToSql()

	likeDests := []*likeDest{}
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

		if _, has := this.postLikersMap[like.PostId]; !has {
			this.postLikersMap[like.PostId] = make(uuid.UUIDs, 0)
		}
		this.postLikersMap[like.PostId] = append(this.postLikersMap[like.PostId], like.UserId)
	}
}

func (this *Presenter) fetchLikeAdditionals() {
	if this.err != nil {
		return
	}

	sql, args, _ := this.qb.
		Select("a.post_id", "count", "b.user_id IS NOT NULL AS liked").
		FromSelect(
			this.qb.
				Select("post_id, count(*)").
				From("post_likes").
				Where(squirrel.Eq{"post_id": this.postIds}).
				GroupBy("post_id"),
			"a",
		).
		JoinClause("LEFT OUTER JOIN post_likes AS b").
		Suffix("ON a.post_id = b.post_id AND b.user_id = ?", this.auth.UserId).
		ToSql()

	additional := []struct {
		PostId uuid.UUID `db:"post_id"`
		Count  int       `db:"count"`
		Liked  bool      `db:"liked"`
	}{}
	err := this.db.Select(&additional, sql, args...)
	if err != nil {
		this.err = errorFrom(err)
		return
	}

	for _, postId := range this.postIds {
		this.out[postId] = &presenterdto.Likes{
			RandomFiveLikers: make([]*presenterdto.User, 0),
		}
	}

	for _, add := range additional {
		this.out[add.PostId].Count = add.Count
		this.out[add.PostId].Liked = add.Liked
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

	this.usersMap = map[uuid.UUID]*presenterdto.User{}
	for _, user := range this.users {
		this.usersMap[user.Id] = user
	}
}

func (this *Presenter) output() (map[uuid.UUID]*presenterdto.Likes, *errors.Error) {
	if this.err != nil {
		err := this.err
		this.reset()
		return nil, err
	}

	for postId, likerIds := range this.postLikersMap {
		for _, userId := range likerIds {
			this.out[postId].RandomFiveLikers = append(this.out[postId].RandomFiveLikers, this.usersMap[userId])
		}
	}

	out := this.out
	this.reset()
	return out, nil
}

func (this *Presenter) reset() {
	this.err = nil
	this.auth = nil
	this.postIds = make(uuid.UUIDs, 0)
	this.userIdsToFetch = make(uuid.UUIDs, 0)
	this.userIdsToFetchMap = make(map[uuid.UUID]struct{})
	this.postLikersMap = make(map[uuid.UUID]uuid.UUIDs)
	this.users = make([]*presenterdto.User, 0)
	this.usersMap = make(map[uuid.UUID]*presenterdto.User)
	this.out = make(map[uuid.UUID]*presenterdto.Likes)
}
