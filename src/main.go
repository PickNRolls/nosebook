package main

import (
	"nosebook/src/application/services/socket"
	roothttp "nosebook/src/deps_root/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"nosebook/src/infra/postgres"
)

func main() {
	db := postgres.Connect()
	hub := socket.NewHub()

	rootHttp := roothttp.New(db, hub)
	rootHttp.Run("0.0.0.0:8080")
}
