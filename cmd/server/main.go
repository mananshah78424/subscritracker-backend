package main

import (
	"context"
	"log"
	"subscritracker/pkg/application"
	"subscritracker/pkg/auth"
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

/*
Register routes function
*/
func registerRoutes(app *application.App) error {
	log.Println("Registering routes!")
	app.Echo.GET("/auth/google/login", auth.GoogleLoginhandler)

	return nil
}
