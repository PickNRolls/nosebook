package rootconvservice

import (
	"context"
	presenterdto "nosebook/src/application/presenters/dto"
	presentermessage "nosebook/src/application/presenters/message"
	domainchat "nosebook/src/domain/chat"
	"nosebook/src/errors"
	"nosebook/src/infra/rabbitmq"

	"github.com/google/uuid"
)

type notifier struct {
	rmqConn   *rabbitmq.Connection
	presenter *presentermessage.Presenter
}

func newNotifier(rmqConn *rabbitmq.Connection, presenter *presentermessage.Presenter) *notifier {
  out := &notifier{
		rmqConn:   rmqConn,
		presenter: presenter,
	}

  return out
}

func (this *notifier) NotifyAbout(userId uuid.UUID, chat *domainchat.Chat) *errors.Error {
	events := chat.Events()

	rmqCh, err := errors.Using(this.rmqConn.Channel())
	defer rmqCh.Close()

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

	for _, event := range events {
		if event.Type() == domainchat.MESSAGE_SENT {
			messageSent := event.(*domainchat.MessageSentEvent)

			messageMap, err := this.presenter.FindByIds(context.TODO(), []uuid.UUID{messageSent.Message.Id})
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

			err = errors.From(rmqCh.Publish(
				"notifications",
				userId.String(),
				false,
				false,
				rabbitmq.Publishing{
					ContentType: "text/json",
					Body:        json,
				}))
			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}
