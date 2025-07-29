package auth

import (
	"log"
	"subscritracker/pkg/account"
	"subscritracker/pkg/application"
	"subscritracker/pkg/models"
	"time"

	"github.com/labstack/echo/v4"
)

func SaveUserToDB(c echo.Context, userInfo map[string]interface{}) error {
	log.Println("Saving user to database with userInfo", userInfo)
	app := c.Get("app").(*application.App)

	accountDetails := &models.Account{
		GoogleID:          userInfo["id"].(string),
		Email:             userInfo["email"].(string),
		Name:              userInfo["name"].(string),
		PictureURL:        userInfo["picture"].(string),
		VerifiedEmail:     userInfo["verified_email"].(bool),
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
		return err
	}

	return nil

}
