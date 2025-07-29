package auth

import (
	"log"
	"subscritracker/pkg/account"
	"subscritracker/pkg/application"
	"subscritracker/pkg/models"
	"subscritracker/pkg/utils"
	"time"

	"github.com/labstack/echo/v4"
)

func SaveGoogleLoggedInUserToDb(c echo.Context, userInfo map[string]interface{}) (map[string]interface{}, error) {
	app := c.Get("app").(*application.App)

	accountDetails := &models.Account{
		GoogleID: func() string {
			id, ok := userInfo["id"].(string)
			if !ok {
				log.Println("Invalid type for 'id' in userInfo")
				return ""
			}
			return id
		}(),
		Email: func() string {
			email, ok := userInfo["email"].(string)
			if !ok {
				log.Println("Invalid type for 'email' in userInfo")
				return ""
			}
			return email
		}(),
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
		VerifiedEmail: func() bool {
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

	err := account.CreateAccount(app, accountDetails)
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
