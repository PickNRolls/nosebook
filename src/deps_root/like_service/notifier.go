package rootlikeservice

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	presenteruser "nosebook/src/application/presenters/user"
	domainlike "nosebook/src/domain/like"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"
	"nosebook/src/infra/rabbitmq"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type notifier struct {
	db            *sqlx.DB
	rmqConn       *rabbitmq.Connection
	userPresenter *presenteruser.Presenter
}

func newNotifier(rmqConn *rabbitmq.Connection, db *sqlx.DB, userPresenter *presenteruser.Presenter) *notifier {
	return &notifier{
		db:            db,
		userPresenter: userPresenter,
		rmqConn:       rmqConn,
	}
}

func (this *notifier) notify(resourceName string, eventType string, userId uuid.UUID, like *domainlike.Like) *errors.Error {
  rmqCh, err := errors.Using(this.rmqConn.Channel())
  defer rmqCh.Close()
	if err != nil {
		return err
	}

	qb := querybuilder.New()
	resourceId := like.Resource.Id()
	payload := struct {
		Id      uuid.UUID          `db:"id" json:"id"`
		Message string             `db:"message" json:"message"`
		Liker   *presenterdto.User `json:"liker"`
	}{}

	sql, args, _ := qb.Select("id", "message").
		From(resourceName).
		Where("id = ?", resourceId).
		ToSql()
	err = errors.From(this.db.Get(&payload, sql, args...))
	if err != nil {
		return err
	}

	userMap, err := this.userPresenter.FindByIds(context.TODO(), []uuid.UUID{like.Owner.Id()})
	if err != nil {
		return err
	}

	payload.Liker = userMap[like.Owner.Id()]

	json, err := (&presenterdto.Event{
		Type:    eventType,
		Payload: payload,
	}).ToJson()
	if err != nil {
		return err
	}

	err = errors.From(rmqCh.ExchangeDeclare(
		"notifications",
		"direct",
		false,
		false,
		false,
		false,
		nil,
	))
	if err != nil {
		return err
	}

	err = errors.From(rmqCh.Publish(
		"notifications",
		userId.String(),
		false,
		false,
		rabbitmq.Publishing{
			ContentType: "text/json",
			Body:        json,
      Expiration: "0",
		}))
	if err != nil {
		return err
	}

	return nil
}

func (this *notifier) NotifyAbout(userId uuid.UUID, like *domainlike.Like) *errors.Error {
	if !like.Value {
		return nil
	}

	switch like.Resource.Type() {
	case domainlike.POST_RESOURCE:
		return this.notify("posts", "post_liked", userId, like)

	case domainlike.COMMENT_RESOURCE:
		return this.notify("comments", "comment_liked", userId, like)
	}

	return nil
}
