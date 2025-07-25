package application

import (
	"context"
	"log"
	"subscritracker/config"
	"subscritracker/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
)

type App struct {
	Config   *config.Config
	Database *bun.DB
	Echo     *echo.Echo
}

func NewApp(ctx context.Context) (*App, error) {
	db, err := utils.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database while creating application: %v", err)
	}

	log.Println("Connected to database while creating application")

	app := &App{
		Config:   config.GetConfig(),
		Database: db,
		Echo:     echo.New(),
	}

	return app, nil
}
