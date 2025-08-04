package auth

import (
	"subscritracker/pkg/application"
	"subscritracker/pkg/utils"
)

func RegisterRoutes(app *application.App) {
	// Public routes (no authentication required)
	app.Echo.POST("/v1/auth/signup", SignUpHandler)
	app.Echo.POST("/v1/auth/login", LoginHandler)
	app.Echo.GET("/v1/auth/verify-email", VerifyEmailHandler)
	app.Echo.POST("/v1/auth/forgot-password", ForgotPasswordHandler)

	// Oauth Rules
	app.Echo.GET("/auth/google/login", GoogleLoginHandler)
	app.Echo.GET("/auth/google/callback", GoogleCallBackHandler)

	// Protected routes (require authentication)
	app.Echo.POST("/v1/auth/logout", LogoutHandler, utils.AuthMiddleware)
	app.Echo.GET("/v1/auth/session", CheckSessionHandler, utils.AuthMiddleware)

}
