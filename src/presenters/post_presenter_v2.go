package presenters

import (
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	"nosebook/src/presenters/post_presenter_v2/dto"
	"nosebook/src/services/auth"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostPresenterV2 struct {
	db sqlx.DB
}

type PostDest struct {
	id        uuid.UUID `db:"id"`
	authorId  uuid.UUID `db:"author_id"`
	ownerId   uuid.UUID `db:"owner_id"`
	message   string    `db:"message"`
	createdAt time.Time `db:"created_at"`
}

func newError(message string) *errors.Error {
	return errors.New("PostQueryError", message)
}

func (p *PostPresenterV2) FindByFilter(filter *dto.FindInputDTO, a *auth.Auth) *dto.FindOutputDTO {
	result := &dto.FindOutputDTO{}

	if filter.OwnerId == uuid.Nil && filter.AuthorId == uuid.Nil {
		result.Err = newError("Отсутствует фильтр")
		return result
	}

	qb := postgres.NewSquirrel()

	query := qb.Select(
		"id", "author_id", "owner_id", "message", "created_at",
	).From(
		"posts",
	).OrderBy(
		"created_at DESC, id DESC",
	).Limit(10)

	if filter.OwnerId != uuid.Nil {
		query = query.Where(
			"owner_id = ?", filter.OwnerId,
		)
	}

	if filter.AuthorId != uuid.Nil {
		query = query.Where(
			"author_id = ?", filter.AuthorId,
		)
	}

	posts := []PostDest{}
	sql, args, _ := query.ToSql()
	err := errors.From(p.db.Select(&posts, sql, args...))
	if err != nil {
		result.Err = err
		return result
	}

	if len(posts) == 0 {
		return result
	}

	userMap := make(map[uuid.UUID]bool)
	for _, post := range posts {
		if _, has := userMap[post.authorId]; !has {
			userMap[post.authorId] = true
		}

		if _, has := userMap[post.ownerId]; !has {
			userMap[post.ownerId] = true
		}
	}
	userIds := make([]uuid.UUID, 1)
	for id := range userMap {
		userIds = append(userIds, id)
	}

	sql, args, _ = qb.Select(
		"id", "first_name", "last_name", "nickname", "created_at",
	).From(
		"users",
	).Where(
		squirrel.Eq{"id": userIds},
	).ToSql()
	users := []dto.UserDTO{}
	err = errors.From(p.db.Select(&users, sql, args...))
	if err != nil {
		result.Err = err
		return result
	}

	for _, post := range posts {
		postDTO := &dto.PostDTO{}
		postDTO.Id = post.id
		postDTO.Message = post.message
		postDTO.CreatedAt = post.createdAt

		for _, user := range users {
			if post.authorId == user.Id {
				postDTO.Author = &user
			}

			if post.ownerId == user.Id {
				postDTO.Owner = &user
			}
		}

		result.Data = append(result.Data, postDTO)
	}

	// result.Next = encodeNext(posts[len(posts)-1])

	return result
}
