package main

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/cmeyer18/weather-common/v4/sql/db_migrations/v4/enviornment"
)

func main() {
	env := enviornment.ProcessEnvironmentVariables()

	connectionDetails := "postgres://" + env.PostgresUsername + ":" + env.PostgresPassword + "@" + env.PostgresAddress + ":" + env.PostgresPort + "/" + env.PostgresDB + "?sslmode=disable"

	db, err := sql.Open("postgres", connectionDetails)
	if err != nil {
		log.Fatal(err)
		return
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
		return
	}
}
