package presentercomment

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	cursorquery "nosebook/src/lib/cursor_query"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Presenter struct {
	db            *sqlx.DB
	likePresenter LikePresenter
	userPresenter UserPresenter
	permissions   Permissions
	tracer        trace.Tracer
}

func New(db *sqlx.DB, likePresenter LikePresenter, userPresenter UserPresenter, permissions Permissions) *Presenter {
	return &Presenter{
		db:            db,
		likePresenter: likePresenter,
		userPresenter: userPresenter,
		permissions:   permissions,
		tracer:        noop.Tracer{},
	}
}

func (this *Presenter) WithTracer(tracer trace.Tracer) *Presenter {
	this.tracer = tracer

	return this
}

func errOut(err error) *FindByFilterOutput {
	return errMsgOut(err.Error())
}

func errMsgOut(message string) *FindByFilterOutput {
	return &FindByFilterOutput{
		Err: newError(message),
	}
}

func (this *Presenter) FindById(parent context.Context, id string, auth *auth.Auth) *comment {
	out := this.FindByFilter(parent, FindByFilterInput{
		Ids: []string{id},
	}, auth)

	if out.Data != nil && len(out.Data) > 0 {
		return out.Data[0]
	}

	return nil
}

func (this *Presenter) FindByFilter(parent context.Context, input FindByFilterInput, auth *auth.Auth) *FindByFilterOutput {
	ctx, span := this.tracer.Start(parent, "comment_presenter.find_by_filter")
	defer span.End()

	var postId uuid.UUID
	if input.PostId != "" {
		var err error
		postId, err = uuid.Parse(input.PostId)
		if err != nil {
			return errOut(err)
		}
	}

	var ids []uuid.UUID
	if input.Ids != nil && len(input.Ids) != 0 {
		ids = make([]uuid.UUID, len(input.Ids))
		for i, id := range input.Ids {
			u, err := uuid.Parse(id)
			if err != nil {
				return errOut(err)
			}

			ids[i] = u
		}
	}

	qb := querybuilder.New()
	query := qb.
		Select("id", "author_id", "message", "created_at").
		From("comments as c").
		Where("removed_at is null").
		Join("post_comments as pc on c.id = pc.comment_id")

	if postId != uuid.Nil {
		query = query.Column("post_id").Where("post_id = ?", postId)
	}

	if ids != nil {
		query = query.Where(squirrel.Eq{"id": ids})
	}

	dest := []*Dest{}
	_, span = this.tracer.Start(ctx, "sql_query")
	cursorQueryOut, error := cursorquery.Do(this.db, &cursorquery.Input[*Dest]{
		Query: query,
		Next:  input.Next,
		Prev:  input.Prev,
		Last:  input.Last,
		Order: &order{},
		Limit: input.Limit,
	}, &dest)
	span.End()
	if error != nil {
		return errOut(error)
	}

  likeCtx, span := this.tracer.Start(ctx, "like_presenter.find_by_comment_ids")
	likesMap, err := this.likePresenter.FindByCommentIds(likeCtx, func() uuid.UUIDs {
		ids := make(uuid.UUIDs, len(dest))
		for i, destComment := range dest {
			ids[i] = destComment.Id
		}
		return ids
	}(), auth)
  span.End()
	if err != nil {
		return errOut(err)
	}

  userCtx, span := this.tracer.Start(ctx, "user_presenter.find_by_ids")
	userMap, err := func() (map[uuid.UUID]*presenterdto.User, *errors.Error) {
		ids := uuid.UUIDs{}
		for _, destComment := range dest {
			ids = append(ids, destComment.AuthorId)
		}

		return this.userPresenter.FindByIds(userCtx, ids)
	}()
  span.End()
	if err != nil {
		errOut(err)
	}

	output := &FindByFilterOutput{
		Data:       make([]*comment, len(dest)),
		TotalCount: cursorQueryOut.TotalCount,
		Next:       cursorQueryOut.Next,
		Prev:       cursorQueryOut.Prev,
	}

	for i, destComment := range dest {
		output.Data[i] = &comment{
			Id:      destComment.Id,
			Author:  userMap[destComment.AuthorId],
			Message: destComment.Message,
			Likes:   likesMap[destComment.Id],
			Permissions: &presenterdto.Permissions{
				Remove: this.permissions.CanRemoveBy(destComment, auth.UserId),
				Update: this.permissions.CanUpdateBy(destComment, auth.UserId),
			},
			CreatedAt: destComment.CreatedAt,
		}
	}

	return output
}

func (this *Presenter) FindByPostId(parent context.Context, id uuid.UUID, auth *auth.Auth) *FindByFilterOutput {
	return this.FindByFilter(parent, FindByFilterInput{
		PostId: id.String(),
		Limit:  5,
	}, auth)
}
