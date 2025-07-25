package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func GoogleLoginhandler(c echo.Context) error {
	url := GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}
