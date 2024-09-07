package rootconvservice

import (
	presenterdto "nosebook/src/application/presenters/dto"
	presentermessage "nosebook/src/application/presenters/message"
	"nosebook/src/application/services/socket"
	domainchat "nosebook/src/domain/chat"
	"nosebook/src/errors"

	"github.com/google/uuid"
)

type socketNotifier struct {
	hub       *socket.Hub
	userId    uuid.UUID
	presenter *presentermessage.Presenter
}

func (this *socketNotifier) Notify(chat *domainchat.Chat) *errors.Error {
	events := chat.Events()

	for _, event := range events {
		if event.Type() == domainchat.MESSAGE_SENT {
			messageSent := event.(*domainchat.MessageSentEvent)

			messageMap, err := this.presenter.FindByIds([]uuid.UUID{messageSent.Message.Id})
			if err != nil {
				return err
			}

			message := messageMap[messageSent.Message.Id]
			json, err := (&presenterdto.Event{
				Type:    "new_message",
				Payload: message,
			}).ToJson()
			if err != nil {
				return err
			}

			this.hub.Broadcast(json, &socket.BroadcastFilter{
				UserId: this.userId,
			})
		}
	}

	return nil
}
