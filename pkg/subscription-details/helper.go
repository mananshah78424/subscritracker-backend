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
			sd.next_due_date,
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
	if filters.NextDueDateFrom != nil {
		query += ` AND sd.next_due_date >= ?`
		args = append(args, *filters.NextDueDateFrom)
	}

	if filters.NextDueDateTo != nil {
		query += ` AND sd.next_due_date <= ?`
		args = append(args, *filters.NextDueDateTo)
	}

	// Add sorting
	if filters.SortBy != "" {
		// Map user-friendly sort field names to actual database column names
		sortFieldMapping := map[string]string{
			"monthly_bill":  "sd.monthly_bill",
			"next_due_date": "sd.next_due_date",
			"start_date":    "sd.start_date",
			"status":        "sd.status",
			"channel_name":  "sc.channel_name",
		}

		sortField := sortFieldMapping[filters.SortBy]
		query += ` ORDER BY ` + sortField

		// Add sort order
		if filters.SortOrder == "desc" {
			query += ` DESC`
		} else {
			query += ` ASC` // Default to ascending
		}
	}

	err := app.Database.NewRaw(query, args...).Scan(context.Background(), &subscriptionDetailsByAccountID)
	if err != nil {
		log.Println("Error getting subscription details: because of database error", err)
		return nil, err
	}

	return subscriptionDetailsByAccountID, nil
}

// CalculateNextDueDate calculates the next due date based on start date, due type, and due day of month
// This function should be used for new subscriptions to determine when the first billing should occur
func CalculateNextDueDate(dueType string, dueDayOfMonth int, startDate time.Time) time.Time {
	currentDate := time.Now()
	// If the due day of month is 1, then simply add a month to the current date and return the 1st of the next month
	if dueDayOfMonth == 1 {
		nextDueDate := currentDate.AddDate(0, 1, 0)
		return time.Date(nextDueDate.Year(), nextDueDate.Month(), 1, 0, 0, 0, 0, nextDueDate.Location())
	}

	// If the due day of month is not 1, then we need to calculate the next due date based on the current date.
	// If current date > due day of month, then add a month to the current date and return the due day of the next month.

	if currentDate.Day() > dueDayOfMonth {
		nextDueDate := currentDate.AddDate(0, 1, 0)
		return time.Date(nextDueDate.Year(), nextDueDate.Month(), dueDayOfMonth, 0, 0, 0, 0, nextDueDate.Location())
	}

	// If current date < due day of month, then return the due day of the current month.
	nextDueDate := time.Date(currentDate.Year(), currentDate.Month(), dueDayOfMonth, 0, 0, 0, 0, currentDate.Location())
	return nextDueDate

}
