package main

import (
	"log"
	"os"
	"strconv"
	"subscritracker/pkg/utils"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	direction := os.Args[1]
	switch direction {
	case "up":
		log.Println("Running up migrations")
	case "down":
		log.Println("Running down migrations")
	case "rollback":
		log.Println("Rolling back latest migration")
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Force command requires a version number")
		}
		log.Printf("Forcing database version to %s", os.Args[2])
	default:
		log.Fatal("Invalid direction, must be up, down, rollback, or force")
	}
	db, databaseErr := utils.NewDatabase()
	if databaseErr != nil {
		log.Fatalf("Failed to setup the database: %v\n", databaseErr)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	switch direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()
	case "rollback":
		err = m.Steps(-1) // Rollback 1 step (latest migration)
	case "force":
		version := os.Args[2]
		versionInt, err := strconv.Atoi(version)
		if err != nil {
			log.Fatalf("Invalid version number: %s", version)
		}
		err = m.Force(versionInt)
	}

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}
