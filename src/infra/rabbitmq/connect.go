package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError[T any](data T, err error) T {
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}

	return data
}

func Connect() (*amqp.Connection, *amqp.Channel) {
	conn := failOnError(amqp.Dial("amqp://guest:guest@rabbitmq:5672/"))
	ch := failOnError(conn.Channel())
	return conn, ch
}
