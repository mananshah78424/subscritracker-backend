package subscriptiondetails

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"subscritracker/pkg/application"
	"subscritracker/pkg/models"

	"github.com/labstack/echo/v4"
)

func CreateSubscriptionDetails(c echo.Context, subscriptionDetails models.Subscription_Details) (models.Subscription_Details, error) {
	app := c.Get("app").(*application.App)

	// Set timestamps
	subscriptionDetails.CreatedAt = time.Now()
	subscriptionDetails.UpdatedAt = time.Now()

	_, err := app.Database.NewInsert().
		Model(&subscriptionDetails).
		Exec(context.Background())
	if err != nil {
		log.Println("Error creating subscription details:", err)
		return models.Subscription_Details{}, err
	}

	return subscriptionDetails, nil
}

func GetSubscriptionDetailsByID(c echo.Context, id int) (models.Subscription_Details, error) {
	app := c.Get("app").(*application.App)

	subscriptionDetails := models.Subscription_Details{}
	err := app.Database.NewSelect().
		Model(&subscriptionDetails).
		Where("id = ?", id).
		Scan(context.Background())
	if err != nil {
		// Check if it's a "no rows" error from Bun ORM
		if err.Error() == "sql: no rows in result set" || err.Error() == "no rows in result set" {
			return models.Subscription_Details{}, errors.New("subscription details not found")
		}
		return models.Subscription_Details{}, err
	}

	return subscriptionDetails, nil
}

func CheckSubscriptionDetailsExists(c echo.Context, id int) (bool, error) {
	_, err := GetSubscriptionDetailsByID(c, id)
	if err != nil {
		// If the error is "not found", return false without error
		if err.Error() == "subscription details not found" {
			return false, nil
		}
		// For other errors, return the error
		return false, err
	}

	return true, nil
}

func CheckExistingSubscriptionByChannel(app *application.App, accountID, channelID int) (bool, error) {
	var subscription models.Subscription_Details
	err := app.Database.NewSelect().
		Model(&subscription).
		Where("account_id = ? AND subscription_channel_id = ?", accountID, channelID).
		Limit(1).
		Scan(context.Background())

	if err != nil {
		// Check if it's a "no rows" error
		if strings.Contains(err.Error(), "no rows in result set") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
