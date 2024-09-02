package rootfriendshipservice

import (
	"nosebook/src/application/services/friendship"
	domainfriendship "nosebook/src/domain/friendship"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db              *sqlx.DB
	requesterId     uuid.UUID
	responderId     uuid.UUID
	onlyAccepted    bool
	onlyNotAccepted bool
}

func newRepository(db *sqlx.DB) friendship.Repository {
	return &repository{
		db: db,
	}
}

func (this *repository) RequesterId(id uuid.UUID) friendship.Repository {
	this.requesterId = id
	return this
}

func (this *repository) ResponderId(id uuid.UUID) friendship.Repository {
	this.responderId = id
	return this
}

func (this *repository) OnlyAccepted() friendship.Repository {
	this.onlyAccepted = true
	this.onlyNotAccepted = false
	return this
}

func (this *repository) OnlyNotAccepted() friendship.Repository {
	this.onlyNotAccepted = true
	this.onlyAccepted = false
	return this
}

func (this *repository) FindOne() *domainfriendship.FriendRequest {
	defer func() {
		this.requesterId = uuid.Nil
		this.responderId = uuid.Nil
		this.onlyNotAccepted = false
		this.onlyAccepted = false
	}()

	qb := querybuilder.New()
	query := qb.
		Select(
			"requester_id",
			"responder_id",
			"message",
			"accepted",
			"viewed",
			"created_at",
		).
		From("friendship_requests").
		Where("requester_id = ?", this.requesterId).
		Where("responder_id = ?", this.responderId)

	if this.onlyAccepted {
		query = query.Where("accepted = true")
	}

	if this.onlyNotAccepted {
		query = query.Where("accepted = false")
	}

	sql, args, _ := query.ToSql()

	request := domainfriendship.FriendRequest{}
	err := this.db.Get(&request, sql, args...)
	if err != nil {
		return nil
	}

	return &request
}

func (this *repository) Save(request *domainfriendship.FriendRequest) *errors.Error {
	qb := querybuilder.New()
	events := request.Events()

	for _, event := range events {
		switch event.Type() {
		case domainfriendship.CREATED:
			sql, args, _ := qb.
				Insert("friendship_requests").
				Columns(
					"requester_id",
					"responder_id",
					"message",
					"accepted",
					"viewed",
					"created_at",
				).
				Values(
					request.RequesterId,
					request.ResponderId,
					request.Message,
					request.Accepted,
					request.Viewed,
					request.CreatedAt,
				).
				ToSql()

			_, err := this.db.Exec(sql, args...)
			if err != nil {
				return errors.From(err)
			}

		case domainfriendship.ACCEPTED:
			fallthrough
		case domainfriendship.DENIED:
			sql, args, _ := qb.
				Update("friendship_requests").
				Set(
					"accepted",
					request.Accepted,
				).
				Set(
					"viewed",
					request.Viewed,
				).
				Where("requester_id = ?", request.RequesterId).
				Where("responder_id = ?", request.ResponderId).
				ToSql()

			_, err := this.db.Exec(sql, args...)
			if err != nil {
				return errors.From(err)
			}

		case domainfriendship.REMOVED:
			event := event.(*domainfriendship.RemovedEvent)
			sql, args, _ := qb.
				Update("friendship_requests").
				Set(
					"accepted",
					request.Accepted,
				).
				Set(
					"responder_id",
					request.ResponderId,
				).
				Set(
					"requester_id",
					request.RequesterId,
				).
				Where("requester_id = ?", event.PreviousRequesterId).
				Where("responder_id = ?", event.PreviousResponderId).
				ToSql()

			_, err := this.db.Exec(sql, args...)
			if err != nil {
				return errors.From(err)
			}

		case domainfriendship.VIEWED:
			sql, args, _ := qb.
				Update("friendship_requests").
				Set(
					"viewed",
					request.Viewed,
				).
				Where("requester_id = ?", request.RequesterId).
				Where("responder_id = ?", request.ResponderId).
				ToSql()

			_, err := this.db.Exec(sql, args...)
			if err != nil {
				return errors.From(err)
			}
		}

	}

	return nil
}
