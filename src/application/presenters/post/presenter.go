package presenterpost

import (
	presenterdto "nosebook/src/application/presenters/dto"
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/lib/cursor_query"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Presenter struct {
	db               *sqlx.DB
	userPresenter    UserPresenter
	commentPresenter CommentPresenter
	likePresenter    LikePresenter
	permissions      Permissions
}

func New(
	db *sqlx.DB,
	userPresenter UserPresenter,
	commentPresenter CommentPresenter,
	likePresenter LikePresenter,
	permissions Permissions,
) *Presenter {
	return &Presenter{
		db:               db,
		userPresenter:    userPresenter,
		commentPresenter: commentPresenter,
		likePresenter:    likePresenter,
		permissions:      permissions,
	}
}

func outErr(err error) *FindByFilterOutput {
	return outMsgErr(err.Error())
}

func outMsgErr(message string) *FindByFilterOutput {
	return &FindByFilterOutput{
		Err: newError(message),
	}
}

func outZero() *FindByFilterOutput {
	return &FindByFilterOutput{
		Data: make([]*Post, 0),
	}
}

func (this *Presenter) FindByFilter(input *FindByFilterInput, a *auth.Auth) *FindByFilterOutput {

	dests, cursorQueryOut, err := func() ([]*Dest, *cursorquery.Output, *errors.Error) {
		out := []*Dest{}

		var ownerId uuid.UUID
		var authorId uuid.UUID
		var ids uuid.UUIDs

		if input.OwnerId != "" {
			var err error
			ownerId, err = uuid.Parse(input.OwnerId)
			if err != nil {
				return nil, nil, errorFrom(err)
			}
		}

		if input.AuthorId != "" {
			var err error
			authorId, err = uuid.Parse(input.AuthorId)
			if err != nil {
				return nil, nil, errorFrom(err)
			}
		}

		if input.Ids != nil && len(input.Ids) != 0 {
			ids = make(uuid.UUIDs, len(input.Ids))
			for i, id := range input.Ids {
				u, err := uuid.Parse(id)
				if err != nil {
					return nil, nil, errorFrom(err)
				}

				ids[i] = u
			}
		}

		if ownerId == uuid.Nil && authorId == uuid.Nil && ids == nil {
			return nil, nil, newError("Отсутствует фильтр")
		}

		qb := querybuilder.New()

		query := qb.
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

		cursorQueryOut, err := cursorquery.Do(this.db, &cursorquery.Input{
			Query:    query,
			Next:     input.Cursor,
			Limit:    10,
			OrderAsc: false,
		}, &out)
		if err != nil {
			return nil, nil, errorFrom(err)
		}

		return out, cursorQueryOut, nil
	}()
	if err != nil {
		return outErr(err)
	}

	postIds := make(uuid.UUIDs, len(dests))
	for i, post := range dests {
		postIds[i] = post.Id
	}

	if len(dests) == 0 {
		return outZero()
	}

	usersMap, err := func() (map[uuid.UUID]*presenterdto.User, *errors.Error) {
		userIds := []uuid.UUID{}
		userIdsMap := map[uuid.UUID]struct{}{}
		out := map[uuid.UUID]*presenterdto.User{}

		for _, dest := range dests {
			if _, has := userIdsMap[dest.AuthorId]; !has {
				userIdsMap[dest.AuthorId] = struct{}{}
			}

			if _, has := userIdsMap[dest.OwnerId]; !has {
				userIdsMap[dest.OwnerId] = struct{}{}
			}
		}

		for id := range userIdsMap {
			userIds = append(userIds, id)
		}

		users, err := this.userPresenter.FindByIds(userIds)
		if err != nil {
			return nil, err
		}

		for _, dto := range users {
			out[dto.Id] = dto
		}

		return out, nil
	}()

	commentsMap := map[uuid.UUID]*presenterdto.FindOut[*presenterdto.Comment]{}
	for _, id := range postIds {
		commentsMap[id] = this.commentPresenter.FindByPostId(id, a)
	}

	postLikesMap, err := this.likePresenter.FindByPostIds(postIds, a)

	posts := func() []*Post {
		out := make([]*Post, 0, len(dests))

		for _, dest := range dests {
			postDTO := &Post{}
			postDTO.Id = dest.Id
			postDTO.Author = usersMap[dest.AuthorId]
			postDTO.Owner = usersMap[dest.OwnerId]
			postDTO.Message = dest.Message
			postDTO.CreatedAt = dest.CreatedAt

			postDTO.Permissions = &presenterdto.Permissions{
				Remove: this.permissions.CanRemoveBy(dest, a.UserId),
				Update: this.permissions.CanUpdateBy(dest, a.UserId),
			}

			postDTO.Likes = postLikesMap[dest.Id]

			postDTO.RecentComments = commentsMap[dest.Id]

			out = append(out, postDTO)
		}

		return out
	}()

	return &FindByFilterOutput{
		Data:       posts,
		Next:       cursorQueryOut.Next,
		TotalCount: cursorQueryOut.TotalCount,
	}
}

func (this *Presenter) FindById(id string, a *auth.Auth) *Post {
	out := this.FindByFilter(&FindByFilterInput{
		Ids: []string{id},
	}, a)

	if out.Data != nil && len(out.Data) != 0 {
		return out.Data[0]
	}

	return nil
}
