package auth

import (
	"log"
	"strings"
	"subscritracker/pkg/account"
	"subscritracker/pkg/application"
	"subscritracker/pkg/models"
	"subscritracker/pkg/utils"
	"time"

	"github.com/labstack/echo/v4"
)

func SaveGoogleLoggedInUserToDb(c echo.Context, userInfo map[string]interface{}) (map[string]interface{}, error) {
	app := c.Get("app").(*application.App)

	// Extract Google ID
	var googleID *string
	if id, ok := userInfo["id"].(string); ok && id != "" {
		googleID = &id
	}

	email := func() string {
		email, ok := userInfo["email"].(string)
		if !ok {
			log.Println("Invalid type for 'email' in userInfo")
			return ""
		}
		return email
	}()

	// Check if user already exists by email
	existingAccount, err := account.GetAccountByEmail(app, email)
	if err != nil && !strings.Contains(err.Error(), "sql: no rows in result set") {
		return nil, err
	}

	if existingAccount != nil {
		// User exists, update with Google ID if not already set
		if existingAccount.GoogleID == nil && googleID != nil {
			existingAccount.GoogleID = googleID
			existingAccount.EmailVerified = true // Google accounts are verified
			err = account.UpdateAccount(app, existingAccount)
			if err != nil {
				return nil, err
			}
		}

		// Generate JWT token for existing user
		token, err := utils.GenerateJWT(existingAccount.ID, existingAccount.Email)
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"token":   token,
			"user":    existingAccount,
			"message": "Login successful",
		}, nil
	}

	// Since the user doesn't exist, create a new account
	accountDetails := &models.Account{
		GoogleID: googleID,
		Email:    email,
		Name: func() string {
			name, ok := userInfo["name"].(string)
			if !ok {
				log.Println("Invalid type for 'name' in userInfo")
				return ""
			}
			return name
		}(),
		PictureURL: func() string {
			picture, ok := userInfo["picture"].(string)
			if !ok {
				log.Println("Invalid type for 'picture' in userInfo")
				return ""
			}
			return picture
		}(),
		EmailVerified: func() bool {
			verifiedEmail, ok := userInfo["verified_email"].(bool)
			if !ok {
				log.Println("Invalid type for 'verified_email' in userInfo")
				return false
			}
			return verifiedEmail
		}(),
		Tier:              "free",
		Status:            "active",
		Features:          map[string]interface{}{},
		SubscriptionCount: 0,
		LastLoginAt:       time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	err = account.CreateAccount(app, accountDetails)
	if err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(accountDetails.ID, accountDetails.Email)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token":   token,
		"user":    accountDetails,
		"message": "Login successful",
	}, nil

}

func CreateSignUpAccountBody(app *application.App, email string, password string, name string, givenName string, familyName string, verificationToken string) (*models.Account, error) {
	accountBody := &models.Account{
		Email:             email,
		PasswordHash:      password,
		Name:              name,
		GivenName:         givenName,
		FamilyName:        familyName,
		VerificationToken: verificationToken,
		EmailVerified:     false, // New accounts are not verified until email verification
		Tier:              "free",
		Status:            "active",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Features:          map[string]interface{}{},
		SubscriptionCount: 0,
		LastLoginAt:       time.Now(),
	}

	err := account.CreateAccount(app, accountBody)
	if err != nil {
		log.Println("Error creating account: ", err)
		return nil, err
	}

	// TODO: Send verification email
	// For now, just return success
	return accountBody, nil
}
