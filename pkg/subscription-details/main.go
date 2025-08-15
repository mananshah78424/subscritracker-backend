package subscriptiondetails

import (
	"subscritracker/pkg/application"
	"subscritracker/pkg/utils"
)

func RegisterRoutes(app *application.App) {
	// app.Echo.GET("/v1/subscription-details", GetAllSubscriptionDetailsHandler, utils.AuthMiddleware)
	app.Echo.GET("/v1/subscription-details/:id", GetSubscriptionDetailsHandler, utils.AuthMiddleware)
	app.Echo.POST("/v1/subscription-details", PostSubscriptionDetailsHandler, utils.AuthMiddleware)
	app.Echo.GET("/v1/user-subscription-details", GetUserSubscriptionDetailsHandler, utils.AuthMiddleware)
}
