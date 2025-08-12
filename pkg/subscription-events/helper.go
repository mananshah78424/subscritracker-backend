package subscriptionevents

import (
	"context"
	"log"
	"subscritracker/pkg/application"
	"subscritracker/pkg/models"
	"time"

	"github.com/labstack/echo/v4"
)

func CreateSubscriptionEvent(c echo.Context, request models.Subscription_Event) (models.Subscription_Event, error) {
	app := c.Get("app").(*application.App)

	// Set timestamps
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	_, err := app.Database.NewInsert().
		Model(&request).
		Exec(context.Background())
	if err != nil {
		log.Println("Error creating subscription event:", err)
		return models.Subscription_Event{}, err
	}

	return request, nil
}
