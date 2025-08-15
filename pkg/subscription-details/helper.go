package subscriptiondetails

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"subscritracker/pkg/application"
	"subscritracker/pkg/models"

	"subscritracker/pkg/validator"

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

func GetSubscriptionDetailsByUserIdWithFilters(app *application.App, accountID int, filters *validator.FilterOptions) ([]SubscriptionDetailsByAccountID, error) {
	var subscriptionDetailsByAccountID []SubscriptionDetailsByAccountID

	// Build the base query
	query := `
		SELECT 
			sd.id,
			sd.account_id,
			sd.subscription_channel_id,
			sc.channel_name as subscription_channel_name,
			sc.channel_image_url,
			sd.start_date,
			sd.due_date,
			sd.status,
			sd.monthly_bill,
			sd.reminder_date,
			sd.reminder_time
		FROM subscription_details sd
		JOIN subscription_channels sc ON sd.subscription_channel_id = sc.id
		WHERE sd.account_id = ?
	`

	args := []interface{}{accountID}

	// Add status filter
	if filters.Status != "" {
		query += ` AND sd.status = ?`
		args = append(args, filters.Status)
	}

	// Add cost range filters
	if filters.MinCost != nil {
		query += ` AND sd.monthly_bill >= ?`
		args = append(args, *filters.MinCost)
	}

	if filters.MaxCost != nil {
		query += ` AND sd.monthly_bill <= ?`
		args = append(args, *filters.MaxCost)
	}

	// Add start date range filters
	if filters.StartDateFrom != nil {
		query += ` AND sd.start_date >= ?`
		args = append(args, *filters.StartDateFrom)
	}

	if filters.StartDateTo != nil {
		query += ` AND sd.start_date <= ?`
		args = append(args, *filters.StartDateTo)
	}

	// Add due date range filters
	if filters.DueDateFrom != nil {
		query += ` AND sd.due_date >= ?`
		args = append(args, *filters.DueDateFrom)
	}

	if filters.DueDateTo != nil {
		query += ` AND sd.due_date <= ?`
		args = append(args, *filters.DueDateTo)
	}

	// Add sorting
	if filters.SortBy != "" {
		// Validate sort field to prevent SQL injection
		validSortFields := map[string]string{
			"monthly_bill": "sd.monthly_bill",
			"due_date":     "sd.due_date",
			"start_date":   "sd.start_date",
			"status":       "sd.status",
			"channel_name": "sc.channel_name",
		}

		if sortField, valid := validSortFields[filters.SortBy]; valid {
			query += ` ORDER BY ` + sortField

			// Add sort order
			if filters.SortOrder == "desc" {
				query += ` DESC`
			} else {
				query += ` ASC` // Default to ascending
			}
		}
	}

	err := app.Database.NewRaw(query, args...).Scan(context.Background(), &subscriptionDetailsByAccountID)
	if err != nil {
		log.Println("Error getting subscription details: because of database error", err)
		return nil, err
	}

	return subscriptionDetailsByAccountID, nil
}
