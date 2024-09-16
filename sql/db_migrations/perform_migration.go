package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/cmeyer18/weather-common/v5/data_structures"
	sql2 "github.com/cmeyer18/weather-common/v5/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/cmeyer18/weather-common/v5/sql/db_migrations/v5/enviornment"
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
	if err != nil {
		log.Fatal(err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	preMigrationVersion, _, err := m.Version()
	if err != nil {
		log.Fatal(err)
		return
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
		return
	}

	postMigrationVersion, _, err := m.Version()
	if err != nil {
		log.Fatal(err)
		return
	}

	if preMigrationVersion == 15 && postMigrationVersion == 16 {
		runMigrationOnUserNotifications(db)
	}
}

func runMigrationOnUserNotifications(db *sql.DB) {
	userNotificationTable := sql2.NewPostgresUserNotificationTable(db)
	locationTable := sql2.NewPostgresLocationTable(db)

	notifications, err := userNotificationTable.SelectAll()
	if err != nil {
		log.Fatal(err)
		return
	}

	println("cdm selectingALL")

	for _, notification := range notifications {
		location := data_structures.Location{
			LocationID:                       notification.NotificationId,
			LocationName:                     notification.LocationName,
			LocationType:                     data_structures.LocationType_UserLocation,
			LocationReferenceID:              notification.UserID,
			ZoneCode:                         notification.ZoneCode,
			CountyCode:                       notification.CountyCode,
			Created:                          notification.CreationTime,
			Latitude:                         notification.Lat,
			Longitude:                        notification.Lng,
			AlertOptions:                     notification.AlertOptions,
			ConvectiveOutlookOptions:         notification.ConvectiveOutlookOptions,
			MesoscaleDiscussionNotifications: notification.MesoscaleDiscussionNotifications,
		}

		err = locationTable.Insert(location)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
