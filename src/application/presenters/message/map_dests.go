package presentermessage

import (
	"context"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

func (this *Presenter) mapDests(ctx context.Context, dests []*dest) (map[uuid.UUID]*message, *errors.Error) {
	userMap, err := func() (map[uuid.UUID]*user, *errors.Error) {
		ids := []uuid.UUID{}
		idMap := make(map[uuid.UUID]struct{})

		for _, dest := range dests {
			if _, has := idMap[dest.AuthorId]; !has {
				idMap[dest.AuthorId] = struct{}{}
				ids = append(ids, dest.AuthorId)
			}
		}

		return this.userPresenter.FindByIds(ctx, ids)
	}()
	if err != nil {
		return nil, err
	}

	out := make(map[uuid.UUID]*message, len(dests))
	for _, dest := range dests {
		out[dest.Id] = &message{
			Id:        dest.Id,
			Author:    userMap[dest.AuthorId],
			Text:      dest.Text,
			ChatId:    dest.ChatId,
			CreatedAt: dest.CreatedAt,
		}

		if dest.ReplyTo.Valid {
			out[dest.Id].ReplyTo = &dest.ReplyTo
		}
	}
	return out, nil
}
