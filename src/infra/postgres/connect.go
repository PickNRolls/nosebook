package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

func Connect() *sqlx.DB {
	postgresPasswordBytes, err := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))
	if err != nil {
		log.Fatalln(err)
	}
	postgresPassword := string(postgresPasswordBytes[:len(postgresPasswordBytes)-1])

	db, err := sqlx.Connect("pgx", fmt.Sprintf("postgres://postgres:%s@db:5432/%s", postgresPassword, os.Getenv("POSTGRES_DB")))
	if err != nil {
		log.Fatalln(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	return db
}
