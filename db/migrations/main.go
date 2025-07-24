package main

import (
	"fmt"
	"log"
	"os"
	"subscritracker/pkg/utils"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	direction := os.Args[1]
	switch direction {
	case "up":
		fmt.Println("Running up migrations")
	case "down":
		fmt.Println("Running down migrations")
	default:
		log.Fatal("Invalid direction, must be up or down")
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
	}

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	fmt.Println("Migrations completed successfully")
}
