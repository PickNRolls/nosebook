package presentercomment

import (
	"nosebook/src/infra/postgres"
	cursorquery "nosebook/src/presenters/cursor_query"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Presenter {
	return &Presenter{
		db: db,
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

func (this *Presenter) FindByFilter(input *FindByFilterInput) *FindByFilterOutput {
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

	output := &FindByFilterOutput{
		Data: make([]*comment, len(dest)),
		Next: cursors.Next,
		Prev: cursors.Prev,
	}

	for i, destComment := range dest {
		output.Data[i] = &comment{
			Id:        destComment.Id,
			Author:    nil,
			Message:   destComment.Message,
			CreatedAt: destComment.CreatedAt,
		}
	}

	return output
}

func (this *Presenter) FindByPostId(id uuid.UUID) *FindByFilterOutput {
	return this.FindByFilter(&FindByFilterInput{
		PostId: id.String(),
		Limit:  5,
	})
}
