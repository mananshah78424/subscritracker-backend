package auth

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"subscritracker/pkg/account"
	"subscritracker/pkg/application"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

// GoogleLoginHandler redirects to Google OAuth
func GoogleLoginHandler(c echo.Context) error {
	config := GoogleOauthConfig
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallbackHandler handles OAuth callback
func GoogleCallBackHandler(c echo.Context) error {
	config := GoogleOauthConfig

	// Get the authorization code from the query params
	code := c.QueryParam("code")
	if code == "" {
		log.Println("No authorization code provided")
		return c.String(http.StatusBadRequest, "No authorization code provided")
	}

	// Exchange the authorization code for an access token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Failed to exchange authorization code for access token:", err)
		return c.String(http.StatusInternalServerError, "Failed to exchange authorization code for access token")
	}

	// Get the user's profile information
	userInfo, err := fetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		log.Println("Failed to fetch user info:", err)
		return c.String(http.StatusInternalServerError, "Failed to fetch user info")
	}

	userDetails, err := SaveGoogleLoggedInUserToDb(c, userInfo)
	if err != nil {
		log.Println("Failed to save user to DB:", err)
		return c.String(http.StatusInternalServerError, "Failed to save user to DB")
	}

	// Return the user info as JSON
	return c.JSON(http.StatusOK, userDetails)
}

// fetchGoogleUserInfo fetches the user info from Google
func fetchGoogleUserInfo(accessToken string) (map[string]interface{}, error) {
	// Create a new request with Authorization header
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with Bearer token
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// LogoutHandler logs out the user
// Todo: Implement this
func LogoutHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})

}

// CheckSessionHandler checks if the user is logged in
func CheckSessionHandler(c echo.Context) error {
	user_id := c.Get("user_id").(int)
	app := c.Get("app").(*application.App)

	account, err := account.GetAccountById(app, user_id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get account")
	}

	return c.JSON(http.StatusOK, account)
}
