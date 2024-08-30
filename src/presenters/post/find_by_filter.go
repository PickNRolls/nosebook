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
	posts []*postDTO

	auth              *auth.Auth
	postDests         []*postDest
	postIds           uuid.UUIDs
	likeDests         []*likeDest
	postLikesCountMap map[uuid.UUID]int
	postLikedMap      map[uuid.UUID]struct{}
	postLikersMap     map[uuid.UUID]uuid.UUIDs
	userIdsToFetchMap map[uuid.UUID]struct{}
	userIdsToFetch    uuid.UUIDs
	users             []*userDTO
	usersMap          map[uuid.UUID]*userDTO
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

func (this *findByFilterQuery) FindByFilter(input *FindByFilterInput, a *auth.Auth) *FindByFilterOutput {
	this.auth = a

	this.fetchPosts(input)
	if len(this.postDests) == 0 {
		return this.output()
	}

	this.defineFetchData()
	this.fetchLikes()
	this.fetchLikeAdditionals()
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
			this.err = errors.From(err)
		}
	}

	if input.AuthorId != "" {
		var err error
		authorId, err = uuid.Parse(input.AuthorId)
		if err != nil {
			this.err = errors.From(err)
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
		this.err = errors.From(err)
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

func (this *findByFilterQuery) fetchLikes() {
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

	this.likeDests = []*likeDest{}
	err := this.db.Select(&this.likeDests, sql, args...)
	if err != nil {
		this.err = errors.From(err)
		return
	}

	for _, like := range this.likeDests {
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

func (this *findByFilterQuery) fetchLikeAdditionals() {
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
		this.err = errors.From(err)
		return
	}

	for _, add := range additional {
		this.postLikesCountMap[add.PostId] = add.Count
		if add.Liked {
			this.postLikedMap[add.PostId] = struct{}{}
		}
	}
}

func (this *findByFilterQuery) fetchUsers() {
	if this.err != nil {
		return
	}

	sql, args, _ := this.qb.Select(
		"id", "first_name", "last_name", "nick", "created_at",
	).From(
		"users",
	).Where(
		squirrel.Eq{"id": this.userIdsToFetch},
	).ToSql()
	error := errors.From(this.db.Select(&this.users, sql, args...))
	if error != nil {
		this.err = error
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

	this.posts = make([]*postDTO, 0, len(this.postDests))

	for _, dest := range this.postDests {
		postDTO := &postDTO{}
		postDTO.Id = dest.Id
		postDTO.Author = this.usersMap[dest.AuthorId]
		postDTO.Owner = this.usersMap[dest.OwnerId]
		postDTO.Message = dest.Message
		postDTO.CreatedAt = dest.CreatedAt
		postDTO.Likes = &likesDTO{
			Count:            this.postLikesCountMap[dest.Id],
			RandomFiveLikers: make([]*userDTO, 0),
		}
		_, postDTO.Likes.Liked = this.postLikedMap[dest.Id]

		for _, userId := range this.postLikersMap[dest.Id] {
			postDTO.Likes.RandomFiveLikers = append(postDTO.Likes.RandomFiveLikers, this.usersMap[userId])
		}

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
		this.err = errors.From(err)
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
		output.Posts = this.posts
	}

	this.reset()

	return output
}

func (this *findByFilterQuery) reset() {
	this.err = nil
	this.next = ""
	this.auth = nil
	this.posts = make([]*postDTO, 0)
	this.postDests = make([]*postDest, 0)
	this.postIds = make(uuid.UUIDs, 0)
	this.likeDests = make([]*likeDest, 0)
	this.postLikesCountMap = make(map[uuid.UUID]int)
	this.postLikedMap = make(map[uuid.UUID]struct{})
	this.postLikersMap = make(map[uuid.UUID]uuid.UUIDs)
	this.userIdsToFetchMap = make(map[uuid.UUID]struct{})
	this.userIdsToFetch = make(uuid.UUIDs, 0)
	this.users = make([]*userDTO, 0)
	this.usersMap = make(map[uuid.UUID]*userDTO)
}
