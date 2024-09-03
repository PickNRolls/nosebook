package presenterfriendship

import (
	"nosebook/src/application/services/auth"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type DescribeRelationInput struct {
	SourceUserId  string
	TargetUserIds []string
}

type DescribeRelationOutput struct {
	FriendIds           []uuid.UUID `json:"friendIds,omitempty"`
	PendingResponderIds []uuid.UUID `json:"pendingResponderIds,omitempty"`
	PendingRequesterIds []uuid.UUID `json:"pendingRequesterIds,omitempty"`
}

func (this *Presenter) DescribeRelation(input *DescribeRelationInput, auth *auth.Auth) (*DescribeRelationOutput, *errors.Error) {
	if len(input.TargetUserIds) > 20 {
		return nil, newError("TargetUserIds не может быть больше 20")
	}

	sourceUserId, err := uuid.Parse(input.SourceUserId)
	if err != nil {
		return nil, errorFrom(err)
	}

	targetUserIds, err := func() (uuid.UUIDs, error) {
		out := make(uuid.UUIDs, len(input.TargetUserIds))

		for i, id := range input.TargetUserIds {
			u, err := uuid.Parse(id)
			if err != nil {
				return nil, err
			}

			out[i] = u
		}

		return out, nil
	}()
	if err != nil {
		return nil, errorFrom(err)
	}

	union := querybuilder.Union(
		squirrel.StatementBuilder.
			Select("requester_id, responder_id, accepted, created_at").
			From("friendship_requests").
			Where("responder_id = ?", sourceUserId).
			Where(squirrel.Eq{"requester_id": targetUserIds}),

		squirrel.StatementBuilder.
			Select("requester_id, responder_id, accepted, created_at").
			From("friendship_requests").
			Where("requester_id = ?", sourceUserId).
			Where(squirrel.Eq{"responder_id": targetUserIds}),
	).PlaceholderFormat(squirrel.Dollar)

	sql, args, _ := union.ToSql()
	dests := make([]*describe_relation_dest, 0)
	err = this.db.Select(&dests, sql, args...)
	if err != nil {
		return nil, errorFrom(err)
	}

	out := &DescribeRelationOutput{}
	for _, dest := range dests {
		if dest.Accepted {
			id := dest.RequesterId
			if id == sourceUserId {
				id = dest.ResponderId
			}

			out.FriendIds = append(out.FriendIds, id)
		} else {
			if dest.RequesterId == sourceUserId {
				out.PendingResponderIds = append(out.PendingResponderIds, dest.ResponderId)
			} else {
				out.PendingRequesterIds = append(out.PendingRequesterIds, dest.RequesterId)
			}
		}
	}
	return out, nil
}
