package presenterpost

import (
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	"nosebook/src/presenters/cursor"
	"nosebook/src/services/auth"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const limit = 10

type findByFilterQuery struct {
	db *sqlx.DB
	qb squirrel.StatementBuilderType

	err   *errors.Error
	next  string
	posts []*post

	auth              *auth.Auth
	userPresenter     userPresenter
	commentPresenter  commentPresenter
	likePresenter     likePresenter
	postLikesMap      map[uuid.UUID]*likes
	postDests         []*postDest
	postIds           uuid.UUIDs
	userIdsToFetchMap map[uuid.UUID]struct{}
	userIdsToFetch    uuid.UUIDs
	users             []*user
	usersMap          map[uuid.UUID]*user
	commentsMap       map[uuid.UUID]*comments
}

func newFindByFilterQuery(db *sqlx.DB) *findByFilterQuery {
	output := &findByFilterQuery{
		db: db,
		qb: postgres.NewSquirrel(),
	}

	output.reset()

	return output
}

type postDest struct {
	Id        uuid.UUID `db:"id"`
	AuthorId  uuid.UUID `db:"author_id"`
	OwnerId   uuid.UUID `db:"owner_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

type likeDest struct {
	PostId uuid.UUID `db:"post_id"`
	UserId uuid.UUID `db:"user_id"`
}

func (this *findByFilterQuery) FindByFilter(
	input *FindByFilterInput,
	a *auth.Auth,
	userPresenter userPresenter,
	commentPresenter commentPresenter,
	likePresenter likePresenter,
) *FindByFilterOutput {
	this.auth = a
	this.userPresenter = userPresenter
	this.commentPresenter = commentPresenter
	this.likePresenter = likePresenter

	this.fetchPosts(input)
	if len(this.postDests) == 0 {
		return this.output()
	}

	this.defineFetchData()
	this.fetchComments()
	this.fetchLikes()
	this.fetchUsers()
	this.mapPosts()
	this.addNextCursor()

	return this.output()
}

func (this *findByFilterQuery) fetchPosts(input *FindByFilterInput) {
	var ownerId uuid.UUID
	var authorId uuid.UUID

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

	if ownerId == uuid.Nil && authorId == uuid.Nil {
		this.err = newError("Отсутствует фильтр")
		return
	}

	query := this.qb.
		Select(
			"id", "author_id", "owner_id", "message", "created_at",
		).
		From(
			"posts",
		).
		Where(
			"removed_at IS NULL",
		).
		OrderBy(
			"created_at DESC, id DESC",
		).
		Limit(limit)

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

	if input.Cursor != "" {
		timestamp, id, err := cursor.Decode(input.Cursor)
		if err != nil {
			this.err = err
			return
		}

		query = query.Where(
			"(created_at, id) < (?, ?)",
			timestamp, id,
		)
	}

	sql, args, _ := query.ToSql()
	err := this.db.Select(&this.postDests, sql, args...)
	if err != nil {
		this.err = errorFrom(err)
	}

	this.postIds = make(uuid.UUIDs, len(this.postDests))
	for i, post := range this.postDests {
		this.postIds[i] = post.Id
	}
}

func (this *findByFilterQuery) defineFetchData() {
	if this.err != nil {
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
	if this.err != nil {
		return
	}

	for _, id := range this.postIds {
		this.commentsMap[id] = this.commentPresenter.FindByPostId(id)
	}
}

func (this *findByFilterQuery) fetchLikes() {
	if this.err != nil {
		return
	}

	out, err := this.likePresenter.FindByPostIds(this.postIds, this.auth)
	if err != nil {
		this.err = err
		return
	}

	this.postLikesMap = out
}

func (this *findByFilterQuery) fetchUsers() {
	if this.err != nil {
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
	if this.err != nil {
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

func (this *findByFilterQuery) addNextCursor() {
	if this.err != nil || len(this.posts) == 0 || len(this.posts) < limit {
		return
	}

	last := this.posts[len(this.posts)-1]

	remainingCount := struct {
		Count int `db:"count"`
	}{}

	sql, args, _ := this.qb.Select("count(*)").
		From("posts").
		Where("removed_at IS NULL").
		Where("(created_at, id) < (?, ?)", last.CreatedAt, last.Id).
		ToSql()
	err := this.db.Get(&remainingCount, sql, args...)
	if err != nil {
		this.err = errorFrom(err)
		return
	}

	if remainingCount.Count > 0 {
		this.next = cursor.Encode(last.CreatedAt, last.Id)
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
	this.err = nil
	this.next = ""
	this.auth = nil
	this.userPresenter = nil
	this.commentPresenter = nil
	this.likePresenter = nil
	this.posts = make([]*post, 0)
	this.postDests = make([]*postDest, 0)
	this.postIds = make(uuid.UUIDs, 0)
	this.userIdsToFetchMap = make(map[uuid.UUID]struct{})
	this.userIdsToFetch = make(uuid.UUIDs, 0)
	this.users = make([]*user, 0)
	this.usersMap = make(map[uuid.UUID]*user)
	this.commentsMap = make(map[uuid.UUID]*comments)
}
