package postgres

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func Connect() *sqlx.DB {
	postgresPasswordBytes, err := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))
	if err != nil {
		panic(err)
	}
	postgresPassword := string(postgresPasswordBytes[:len(postgresPasswordBytes)-1])

	db, err := sqlx.Connect("pgx", fmt.Sprintf("postgres://postgres:%s@db:5432/%s", postgresPassword, os.Getenv("POSTGRES_DB")))
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(0)

	return db
}
