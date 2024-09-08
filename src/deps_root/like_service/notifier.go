package rootlikeservice

import (
	presenterdto "nosebook/src/application/presenters/dto"
	presenteruser "nosebook/src/application/presenters/user"
	"nosebook/src/application/services/socket"
	domainlike "nosebook/src/domain/like"
	"nosebook/src/errors"
	querybuilder "nosebook/src/infra/query_builder"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type socketNotifier struct {
	db            *sqlx.DB
	hub           *socket.Hub
	userId        uuid.UUID
	userPresenter *presenteruser.Presenter
}

func (this *socketNotifier) notify(table string, eventType string, like *domainlike.Like) *errors.Error {
	qb := querybuilder.New()
	resourceId := like.Resource.Id()
	payload := struct {
		Id      uuid.UUID          `db:"id" json:"id"`
		Message string             `db:"message" json:"message"`
		Liker   *presenterdto.User `json:"liker"`
	}{}

	sql, args, _ := qb.Select("id", "message").
		From(table).
		Where("id = ?", resourceId).
		ToSql()
	err := errors.From(this.db.Get(&payload, sql, args...))
	if err != nil {
		return err
	}

	userMap, err := this.userPresenter.FindByIds([]uuid.UUID{like.Owner.Id()})
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

	this.hub.Broadcast(json, &socket.BroadcastFilter{
		UserId: this.userId,
	})

	return nil
}

func (this *socketNotifier) NotifyAbout(like *domainlike.Like) *errors.Error {
	if !like.Value {
		return nil
	}

	switch like.Resource.Type() {
	case domainlike.POST_RESOURCE:
		return this.notify("posts", "post_liked", like)

	case domainlike.COMMENT_RESOURCE:
		return this.notify("comments", "comment_liked", like)
	}

	return nil
}
