package postgres

import (
	"fmt"
	"nosebook/src/lib/config"
	"nosebook/src/lib/secret"

	"github.com/jmoiron/sqlx"
)

func Connect() *sqlx.DB {
	db, err := sqlx.Connect("pgx", fmt.Sprintf("postgres://postgres:%s@db:5432/%s", secret.DBPassword, config.DBName))
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
