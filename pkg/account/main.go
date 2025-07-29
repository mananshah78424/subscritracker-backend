package account

import (
	"subscritracker/pkg/application"
)

func RegisterRoutes(app *application.App) {
	app.Echo.GET("/v1/account/:id", GetAccountByIdHandler)
	app.Echo.GET("/v1/account", GetAccountHandler)
	app.Echo.PUT("/v1/account", UpdateAccountHandler)
	app.Echo.GET("/v1/account/stats", GetAccountStatsHandler)
	app.Echo.POST("/v1/account", CreateAccountHandler)
}
