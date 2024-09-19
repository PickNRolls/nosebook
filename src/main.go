package main

import (
	roothttp "nosebook/src/deps_root/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"nosebook/src/infra/postgres"
	"nosebook/src/infra/rabbitmq"
)

func main() {
	db := postgres.Connect()
	rmqConn:= rabbitmq.Connect()
	defer rmqConn.Close()

	rootHttp := roothttp.New(db, rmqConn)
	rootHttp.Run("0.0.0.0:8080")
}
