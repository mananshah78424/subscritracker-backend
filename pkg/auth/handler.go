package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"subscritracker/pkg/account"
	"subscritracker/pkg/application"
	"subscritracker/pkg/utils"
	"subscritracker/pkg/validator"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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
	app := c.Get("app").(*application.App)
	frontendURL := app.Config.Frontend.URL

	// Get the authorization code from the query params
	code := c.QueryParam("code")
	if code == "" {
		log.Println("No authorization code provided")
		// Redirect to frontend with error
		log.Println("Redirecting to login page with error - No authorization code provided")
		return c.Redirect(http.StatusTemporaryRedirect, frontendURL)
	}

	// Exchange the authorization code for an access token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Failed to exchange authorization code for access token:", err)
		log.Println("Redirecting to login page with error %s- Failed to exchange authorization code for access token", err)
		return c.Redirect(http.StatusTemporaryRedirect, frontendURL)
	}

	// Get the user's profile information
	userInfo, err := fetchGoogleUserInfo(token.AccessToken)
	if err != nil {
		log.Println("Failed to fetch user info:", err)
		log.Println("Redirecting to login page with error %s- Failed to fetch user info", err)
		return c.Redirect(http.StatusTemporaryRedirect, frontendURL)
	}

	userDetails, err := SaveGoogleLoggedInUserToDb(c, userInfo)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			// User already exists, try to get the existing user and login
			existingUser, err := account.GetAccountByEmail(app, userInfo["email"].(string))
			if err != nil {
				log.Printf("Failed to get existing user: %v", err)
				log.Println("Redirecting to login page with error %s- Failed to get existing user", err)
				return c.Redirect(http.StatusTemporaryRedirect, frontendURL)
			}

			// Generate JWT token for existing user
			token, err := utils.GenerateJWT(existingUser.ID, existingUser.Email)
			if err != nil {
				log.Printf("Redirecting to login page with error %v- Failed to generate token for existing user", err)
				return c.Redirect(http.StatusTemporaryRedirect, frontendURL)
			}

			// Redirect to frontend with token and user data
			userJSON, err := json.Marshal(existingUser)
			if err != nil {
				log.Println("Redirecting to login page with error %v- Failed to marshal user data", err)
				return c.Redirect(http.StatusTemporaryRedirect, frontendURL)
			}
			redirectURL := fmt.Sprintf("%s/home?token=%s&user=%s",
				frontendURL,
				token,
				url.QueryEscape(string(userJSON)))

			return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		} else {
			log.Printf("Got error here, redirecting to login page", err)
			// Other error, send to login page without error in url
			return c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/login")
		}
	}

	// Redirect to frontend with token and user data
	userJSON, _ := json.Marshal(userDetails["user"])
	redirectURL := fmt.Sprintf("%s/home?token=%s&user=%s",
		frontendURL,
		userDetails["token"],
		url.QueryEscape(string(userJSON)))

	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
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

func SignUpHandler(c echo.Context) error {
	var req validator.SignUpRequest

	// Try to bind the request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
	}

	// Validate required fields
	if err := validator.ValidateSignUp(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	app := c.Get("app").(*application.App)

	// Check if email already exists
	existingAccount, err := account.GetAccountByEmail(app, req.Email)
	if err != nil {
		// Check if it's a "no rows" error (which means email doesn't exist - good for signup)
		if err.Error() == "sql: no rows in result set" {
			// This is expected - email doesn't exist, so we can proceed with signup
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check if email exists"})
		}
	} else if existingAccount != nil {
		// Account exists with this email
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email already exists"})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	// Generate verification token
	verificationToken, err := account.GenerateToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate verification token"})
	}

	accountBody, err := CreateSignUpAccountBody(app, req.Email, string(hashedPassword), req.Name, req.GivenName, req.FamilyName, verificationToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create account"})
	}

	// TODO: Send verification email
	// For now, just return success
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Account created successfully. Please check your email to verify your account.",
		"user_id": accountBody.ID,
	})
}

func LoginHandler(c echo.Context) error {
	var req validator.LoginRequest

	// Try to bind the request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body: " + err.Error()})
	}

	// Validate required fields
	if err := validator.ValidateLogin(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	app := c.Get("app").(*application.App)

	// Get account by email
	accountDetails, err := account.GetAccountByEmail(app, req.Email)
	if err != nil || accountDetails == nil || !accountDetails.EmailVerified || accountDetails.Status != "active" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}
	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(accountDetails.PasswordHash), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	// Update last login
	accountDetails.LastLoginAt = time.Now()
	accountDetails.Status = "active"
	err = account.UpdateAccount(app, accountDetails)
	if err != nil {
		log.Println("Failed to update last login:", err)
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(accountDetails.ID, accountDetails.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":   token,
		"user":    accountDetails,
		"message": "Login successful",
	})
}

// VerifyEmailHandler verifies the email of the user
func VerifyEmailHandler(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Token is required"})
	}

	app := c.Get("app").(*application.App)
	accountDetails, err := account.GetAccountByVerificationToken(app, token)
	if err != nil || accountDetails == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid verification token"})
	}

	accountDetails.EmailVerified = true
	accountDetails.Status = "active"
	accountDetails.VerificationToken = ""
	err = account.UpdateAccount(app, accountDetails)
	if err != nil {
		log.Println("Failed to update account:", err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Email verified successfully. You can now log in.",
	})
}

// ForgotPasswordHandler sends a reset password email to the user
func ForgotPasswordHandler(c echo.Context) error {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	app := c.Get("app").(*application.App)
	accountDetails, err := account.GetAccountByEmail(app, req.Email)
	if err != nil || accountDetails == nil {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "If an account with this email exists, a password reset link has been sent.",
		})
	}

	// Generate reset token
	resetToken, err := account.GenerateToken()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate reset token"})
	}

	// Set reset token and expiration
	accountDetails.ResetToken = resetToken
	accountDetails.ResetTokenExpires = time.Now().Add(time.Hour * 24)
	err = account.UpdateAccount(app, accountDetails)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to set reset token"})
	}

	// TODO: Send password reset email
	// For now, just return success
	return c.JSON(http.StatusOK, map[string]string{
		"message": "If an account with this email exists, a password reset link has been sent.",
	})
}
