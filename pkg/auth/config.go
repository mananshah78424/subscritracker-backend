package auth

import (
	"subscritracker/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig = &oauth2.Config{

	ClientID:     config.GetConfig().GoogleAuth.ClientID,
	ClientSecret: config.GetConfig().GoogleAuth.ClientSecret,
	RedirectURL:  config.GetConfig().GoogleAuth.RedirectURL,
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}
