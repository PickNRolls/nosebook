package rootconvservice

import (
	"nosebook/src/application/services/conversation"
	"nosebook/src/application/services/socket"
	domainchat "nosebook/src/domain/chat"

	"github.com/google/uuid"
)

type socketNotifier struct {
	client socket.Client
}

func (this *socketNotifier) Notify(chat *domainchat.Chat) {
	events := chat.Events()

	for _, event := range events {
		if event.Type() == domainchat.MESSAGE_SENT {
			messageSent := event.(*domainchat.MessageSentEvent)

			this.client.Send() <- []byte(messageSent.Message.Text)
		}
	}
}

type notifierRepository struct {
	hub *socket.Hub
}

func newNotifierRepository(hub *socket.Hub) *notifierRepository {
	return &notifierRepository{
		hub: hub,
	}
}

func (this *notifierRepository) FindByRecipientId(id uuid.UUID) conversation.Notifier {
	client := this.hub.Client(id)
	if client == nil {
		return nil
	}

	return &socketNotifier{
		client: client,
	}
}
