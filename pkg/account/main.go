package account

import (
	"subscritracker/pkg/application"
	"subscritracker/pkg/utils"
)

func RegisterRoutes(app *application.App) {
	app.Echo.GET("/v1/account/:id", GetAccountByIdHandler, utils.AuthMiddleware)
	app.Echo.GET("/v1/account", GetAccountHandler, utils.AuthMiddleware)
	app.Echo.PUT("/v1/account", UpdateAccountHandler, utils.AuthMiddleware)
	app.Echo.GET("/v1/account/stats", GetAccountStatsHandler, utils.AuthMiddleware)
	app.Echo.POST("/v1/account", CreateAccountHandler)
}
