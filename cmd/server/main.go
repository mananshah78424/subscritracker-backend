package main

import (
	"context"
	"log"
	"subscritracker/pkg/account"
	analysis "subscritracker/pkg/analysis"
	"subscritracker/pkg/application"
	"subscritracker/pkg/auth"
	subscription_channels "subscritracker/pkg/subscription-channels"
	subscription_details "subscritracker/pkg/subscription-details"
	subscription_events "subscritracker/pkg/subscription-events"

	"github.com/labstack/echo/v4"
)

/*
Main function
*/
func main() {
	// Create application context
	ctx := context.Background()

	// Create application
	app, err := application.NewApp(ctx)
	if err != nil {
		log.Fatalf("Failed to create application: %v", err)
	}

	// Register routes
	err = registerRoutes(app)
	if err != nil {
		log.Fatalf("Failed to register routes: %v", err)
	}

	// Start server
	if err := app.Echo.Start(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func appMiddleware(app *application.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	}
}

/*
Register routes function
*/
func registerRoutes(app *application.App) error {
	log.Println("Registering routes!")
	app.Echo.Use(appMiddleware(app))
	auth.RegisterRoutes(app)
	account.RegisterRoutes(app)
	subscription_channels.RegisterRoutes(app)
	subscription_details.RegisterRoutes(app)
	subscription_events.RegisterRoutes(app)
	analysis.RegisterRoutes(app)

	return nil
}
