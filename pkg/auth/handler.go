package auth

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func GoogleLoginHandler(c echo.Context) error {
	url := GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

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

	err = SaveUserToDB(c, userInfo)
	if err != nil {
		log.Println("Failed to save user to DB:", err)
		return c.String(http.StatusInternalServerError, "Failed to save user to DB")
	}

	// Return the user info as JSON
	return c.JSON(http.StatusOK, userInfo)
}

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
