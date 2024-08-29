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

type Presenter struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Presenter {
	return &Presenter{
		db: db,
	}
}

func (this *Presenter) FindByFilter(input *FindByFilterInput, a *auth.Auth) *FindByFilterOutput {
	output := &FindByFilterOutput{}
	var ownerId uuid.UUID
	var authorId uuid.UUID

	if input.OwnerId != "" {
		var err error
		ownerId, err = uuid.Parse(input.OwnerId)
		if err != nil {
			output.Err = errors.From(err)
		}
	}

	if input.AuthorId != "" {
		var err error
		authorId, err = uuid.Parse(input.AuthorId)
		if err != nil {
			output.Err = errors.From(err)
		}
	}

	if ownerId == uuid.Nil && authorId == uuid.Nil {
		output.Err = newError("Отсутствует фильтр")
		return output
	}

	qb := postgres.NewSquirrel()

	query := qb.
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
			output.Err = err
			return output
		}

		query = query.Where(
			"(created_at, id) < (?, ?)",
			timestamp, id,
		)
	}

	sql, args, _ := query.ToSql()
	posts := []struct {
		Id        uuid.UUID `db:"id"`
		AuthorId  uuid.UUID `db:"author_id"`
		OwnerId   uuid.UUID `db:"owner_id"`
		Message   string    `db:"message"`
		CreatedAt time.Time `db:"created_at"`
	}{}
	err := this.db.Select(&posts, sql, args...)
	if err != nil {
		output.Err = errors.From(err)
		return output
	}

	if len(posts) == 0 {
		output.Posts = make([]*PostDTO, 0)
		return output
	}

	userMap := make(map[uuid.UUID]struct{})
	for _, post := range posts {
		if _, has := userMap[post.AuthorId]; !has {
			userMap[post.AuthorId] = struct{}{}
		}

		if _, has := userMap[post.OwnerId]; !has {
			userMap[post.OwnerId] = struct{}{}
		}
	}
	userIds := make([]uuid.UUID, 1)
	for id := range userMap {
		userIds = append(userIds, id)
	}

	sql, args, _ = qb.Select(
		"id", "first_name", "last_name", "nick", "created_at",
	).From(
		"users",
	).Where(
		squirrel.Eq{"id": userIds},
	).ToSql()
	users := []*UserDTO{}
	error := errors.From(this.db.Select(&users, sql, args...))
	if error != nil {
		output.Err = error
		return output
	}

	for _, post := range posts {
		postDTO := &PostDTO{}
		postDTO.Id = post.Id
		postDTO.Message = post.Message
		postDTO.CreatedAt = post.CreatedAt

		for _, user := range users {
			if post.AuthorId == user.Id {
				postDTO.Author = user
			}

			if post.OwnerId == user.Id {
				postDTO.Owner = user
			}
		}

		output.Posts = append(output.Posts, postDTO)
	}

	last := output.Posts[len(output.Posts)-1]

	if len(output.Posts) == limit {
		remainingCount := struct {
			Count int `db:"count"`
		}{}

		sql, args, _ = qb.Select("count(*)").
			From("posts").
			Where("removed_at IS NULL").
			Where("(created_at, id) < (?, ?)", last.CreatedAt, last.Id).
			ToSql()
		err := this.db.Get(&remainingCount, sql, args...)
		if err != nil {
			output.Err = errors.From(err)
			output.Posts = nil
			return output
		}

		if remainingCount.Count > 0 {
			output.Next = cursor.Encode(last.CreatedAt, last.Id)
		}
	}

	return output
}
