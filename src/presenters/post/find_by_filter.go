package presenterpost

import (
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	cursorquery "nosebook/src/lib/cursor_query"
	"nosebook/src/services/auth"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const limit = 10

type findByFilterQuery struct {
	db               *sqlx.DB
	qb               squirrel.StatementBuilderType
	auth             *auth.Auth
	userPresenter    userPresenter
	commentPresenter commentPresenter
	likePresenter    likePresenter

	doSkip bool
	err    *errors.Error
	next   string
	posts  []*post

	postLikesMap      map[uuid.UUID]*likes
	postDests         []*dest
	postIds           uuid.UUIDs
	userIdsToFetchMap map[uuid.UUID]struct{}
	userIdsToFetch    uuid.UUIDs
	users             []*user
	usersMap          map[uuid.UUID]*user
	commentsMap       map[uuid.UUID]*comments
}

func newFindByFilterQuery(
	db *sqlx.DB,
	userPresenter userPresenter,
	commentPresenter commentPresenter,
	likePresenter likePresenter,
) *findByFilterQuery {
	output := &findByFilterQuery{
		db:               db,
		qb:               postgres.NewSquirrel(),
		userPresenter:    userPresenter,
		commentPresenter: commentPresenter,
		likePresenter:    likePresenter,
	}

	output.reset()

	return output
}

type likeDest struct {
	PostId uuid.UUID `db:"post_id"`
	UserId uuid.UUID `db:"user_id"`
}

func (this *findByFilterQuery) skip() bool {
	return this.doSkip || this.err != nil
}

func (this *findByFilterQuery) FindByFilter(
	input *FindByFilterInput,
	a *auth.Auth,
) *FindByFilterOutput {
	this.auth = a

	this.fetchPosts(input)
	this.defineFetchData()
	this.fetchComments()
	this.fetchLikes()
	this.fetchUsers()
	this.mapPosts()

	return this.output()
}

func (this *findByFilterQuery) fetchPosts(input *FindByFilterInput) {
	var ownerId uuid.UUID
	var authorId uuid.UUID
	var ids uuid.UUIDs

	if input.OwnerId != "" {
		var err error
		ownerId, err = uuid.Parse(input.OwnerId)
		if err != nil {
			this.err = errorFrom(err)
		}
	}

	if input.AuthorId != "" {
		var err error
		authorId, err = uuid.Parse(input.AuthorId)
		if err != nil {
			this.err = errorFrom(err)
		}
	}

	if input.Ids != nil && len(input.Ids) != 0 {
		ids = make(uuid.UUIDs, len(input.Ids))
		for i, id := range input.Ids {
			u, err := uuid.Parse(id)
			if err != nil {
				this.err = errorFrom(err)
				break
			}

			ids[i] = u
		}
	}

	if ownerId == uuid.Nil && authorId == uuid.Nil && ids == nil {
		this.err = newError("Отсутствует фильтр")
		return
	}

	query := this.qb.
		Select("id", "author_id", "owner_id", "message", "created_at").
		From("posts").
		Where("removed_at is null")

	if ownerId != uuid.Nil {
		query = query.Where(
			"owner_id = ?", ownerId,
		)
	}

	if authorId != uuid.Nil {
		query = query.Where(
			"author_id = ?", authorId,
		)
	}

	if ids != nil {
		query = query.Where(
			squirrel.Eq{"id": ids},
		)
	}

	cursors, err := cursorquery.Do(this.db, &cursorquery.Input{
		Query:    query,
		Next:     input.Cursor,
		Limit:    limit,
		OrderAsc: false,
	}, &this.postDests)
	if err != nil {
		this.err = errorFrom(err)
	}

	this.next = cursors.Next

	this.postIds = make(uuid.UUIDs, len(this.postDests))
	for i, post := range this.postDests {
		this.postIds[i] = post.Id
	}

	if len(this.postDests) == 0 {
		this.doSkip = true
	}
}

func (this *findByFilterQuery) defineFetchData() {
	if this.skip() {
		return
	}

	for _, post := range this.postDests {
		if _, has := this.userIdsToFetchMap[post.AuthorId]; !has {
			this.userIdsToFetchMap[post.AuthorId] = struct{}{}
		}

		if _, has := this.userIdsToFetchMap[post.OwnerId]; !has {
			this.userIdsToFetchMap[post.OwnerId] = struct{}{}
		}
	}
	for id := range this.userIdsToFetchMap {
		this.userIdsToFetch = append(this.userIdsToFetch, id)
	}
}

func (this *findByFilterQuery) fetchComments() {
	if this.skip() {
		return
	}

	for _, id := range this.postIds {
		this.commentsMap[id] = this.commentPresenter.FindByPostId(id, this.auth)
	}
}

func (this *findByFilterQuery) fetchLikes() {
	if this.skip() {
		return
	}

	this.postLikesMap, this.err = this.likePresenter.FindByPostIds(this.postIds, this.auth)
}

func (this *findByFilterQuery) fetchUsers() {
	if this.skip() {
		return
	}

	this.users, this.err = this.userPresenter.FindByIds(this.userIdsToFetch)
	if this.err != nil {
		return
	}

	for _, dto := range this.users {
		this.usersMap[dto.Id] = dto
	}
}

func (this *findByFilterQuery) mapPosts() {
	if this.skip() {
		return
	}

	this.posts = make([]*post, 0, len(this.postDests))

	for _, dest := range this.postDests {
		postDTO := &post{}
		postDTO.Id = dest.Id
		postDTO.Author = this.usersMap[dest.AuthorId]
		postDTO.Owner = this.usersMap[dest.OwnerId]
		postDTO.Message = dest.Message
		postDTO.CreatedAt = dest.CreatedAt

		postDTO.Likes = this.postLikesMap[dest.Id]

		postDTO.RecentComments = this.commentsMap[dest.Id]

		this.posts = append(this.posts, postDTO)
	}
}

func (this *findByFilterQuery) output() *FindByFilterOutput {
	output := &FindByFilterOutput{
		Err:  this.err,
		Next: this.next,
	}

	if output.Err == nil {
		output.Data = this.posts
	}

	this.reset()

	return output
}

func (this *findByFilterQuery) reset() {
	this.doSkip = false
	this.err = nil
	this.next = ""
	this.auth = nil
	this.posts = make([]*post, 0)
	this.postDests = make([]*dest, 0)
	this.postIds = make(uuid.UUIDs, 0)
	this.userIdsToFetchMap = make(map[uuid.UUID]struct{})
	this.userIdsToFetch = make(uuid.UUIDs, 0)
	this.users = make([]*user, 0)
	this.usersMap = make(map[uuid.UUID]*user)
	this.commentsMap = make(map[uuid.UUID]*comments)
}
