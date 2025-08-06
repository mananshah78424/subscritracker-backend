package subscription_channels

import (
	"subscritracker/pkg/application"
	"subscritracker/pkg/utils"
)

func RegisterRoutes(app *application.App) {
	app.Echo.GET("/v1/subscription-channels", GetAllSubscriptionChannelsHandler, utils.AuthMiddleware)
	app.Echo.GET("/v1/subscription-channels/:id", GetSubscriptionChannelsHandler, utils.AuthMiddleware)
	app.Echo.POST("/v1/subscription-channels", PostSubscriptionChannelsHandler, utils.AuthMiddleware)
}
