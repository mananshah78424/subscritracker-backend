package subscriptionevents

import (
	"subscritracker/pkg/application"
	"subscritracker/pkg/utils"
)

func RegisterRoutes(app *application.App) {
	// app.Echo.GET("/v1/subscription-events", GetAllSubscriptionEventsHandler, utils.AuthMiddleware)
	// app.Echo.GET("/v1/subscription-events/:id", GetSubscriptionEventsHandler, utils.AuthMiddleware)
	app.Echo.POST("/v1/subscription-events", PostSubscriptionEventsHandler, utils.AuthMiddleware)
}
