package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	roothttp "nosebook/src/deps_root/http"

	"nosebook/src/infra/postgres"
)

func main() {
	db := postgres.Connect()
	rootHttp := roothttp.New(db)
	rootHttp.Run("0.0.0.0:8080")
}
