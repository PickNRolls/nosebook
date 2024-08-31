package presentercomment

import (
	"nosebook/src/errors"
	"nosebook/src/infra/postgres"
	"nosebook/src/lib/cursor_query"
	presenterdto "nosebook/src/presenters/dto"
	"nosebook/src/services/auth"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	db            *sqlx.DB
	likePresenter likePresenter
	userPresenter userPresenter
}

func New(db *sqlx.DB, likePresenter likePresenter, userPresenter userPresenter) *Presenter {
	return &Presenter{
		db:            db,
		likePresenter: likePresenter,
		userPresenter: userPresenter,
	}
}

func errOut(err error) *FindByFilterOutput {
	return errMsgOut(err.Error())
}

func errMsgOut(message string) *FindByFilterOutput {
	return &FindByFilterOutput{
		Err: newError(message),
	}
}

func (this *Presenter) FindByFilter(input *FindByFilterInput, auth *auth.Auth) *FindByFilterOutput {
	if input.PostId == "" {
		return errMsgOut("Отсутствует фильтр по PostId")
	}

	postId, err := uuid.Parse(input.PostId)
	if err != nil {
		return errOut(err)
	}

	qb := postgres.NewSquirrel()
	query := qb.
		Select("id", "post_id", "author_id", "message", "created_at").
		From("comments as c").
		Where("removed_at is null").
		Where("post_id = ?", postId).
		Join("post_comments as pc on c.id = pc.comment_id")

	dest := []*commentDest{}
	cursors, error := cursorquery.Do(this.db, &cursorquery.Input{
		Query:    query,
		Next:     input.Next,
		Prev:     input.Prev,
		Last:     input.Last,
		OrderAsc: true,
		Limit:    input.Limit,
	}, &dest)
	if error != nil {
		return errOut(error)
	}

	likesMap, err := this.likePresenter.FindByCommentIds(func() uuid.UUIDs {
		ids := make(uuid.UUIDs, len(dest))
		for i, destComment := range dest {
			ids[i] = destComment.Id
		}
		return ids
	}(), auth)

	userMap, err := func() (map[uuid.UUID]*presenterdto.User, *errors.Error) {
		ids := uuid.UUIDs{}
		for _, destComment := range dest {
			ids = append(ids, destComment.AuthorId)
		}

		users, err := this.userPresenter.FindByIds(ids)
		if err != nil {
			return nil, err
		}

		m := map[uuid.UUID]*presenterdto.User{}
		for _, user := range users {
			m[user.Id] = user
		}
		return m, nil
	}()

	output := &FindByFilterOutput{
		Data: make([]*comment, len(dest)),
		Next: cursors.Next,
		Prev: cursors.Prev,
	}

	for i, destComment := range dest {
		output.Data[i] = &comment{
			Id:        destComment.Id,
			Author:    userMap[destComment.AuthorId],
			Message:   destComment.Message,
			Likes:     likesMap[destComment.Id],
			CreatedAt: destComment.CreatedAt,
		}
	}

	return output
}

func (this *Presenter) FindByPostId(id uuid.UUID, auth *auth.Auth) *FindByFilterOutput {
	return this.FindByFilter(&FindByFilterInput{
		PostId: id.String(),
		Limit:  5,
	}, auth)
}
