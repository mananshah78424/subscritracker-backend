package application

import (
	"context"
	"log"
	"subscritracker/config"
	"subscritracker/pkg/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	// Add CORS middleware globally
	app.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{app.Config.Frontend.URL, "http://127.0.0.1:3000"}, // Keep 127.0.0.1 for local development
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	return app, nil
}
