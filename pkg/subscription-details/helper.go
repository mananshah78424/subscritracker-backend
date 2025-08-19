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

// Calculate NextDueDate based on DueType and DueDayOfMonth and Current Date
func CalculateNextDueDate(dueType string, dueDayOfMonth int, currentDate time.Time) time.Time {
	// First, try to create a date in the current month with the specified due day
	nextDueDate := time.Date(currentDate.Year(), currentDate.Month(), dueDayOfMonth, 0, 0, 0, 0, currentDate.Location())

	// If the due day has already passed this month, move to next month
	if nextDueDate.Before(currentDate) || nextDueDate.Equal(currentDate) {
		if dueType == "monthly" {
			nextDueDate = nextDueDate.AddDate(0, 1, 0)
		} else if dueType == "yearly" {
			nextDueDate = nextDueDate.AddDate(1, 0, 0)
		} else if dueType == "weekly" {
			// For weekly, find the next occurrence
			for nextDueDate.Before(currentDate) || nextDueDate.Equal(currentDate) {
				nextDueDate = nextDueDate.AddDate(0, 0, 7)
			}
		} else if dueType == "daily" {
			// For daily, use tomorrow
			nextDueDate = currentDate.AddDate(0, 0, 1)
		}
	}

	// Handle edge case where the due day doesn't exist in the target month
	// (e.g., due day is 31st but target month only has 30 days)
	if dueType == "monthly" || dueType == "yearly" {
		for {
			// Try to create the date
			testDate := time.Date(nextDueDate.Year(), nextDueDate.Month(), dueDayOfMonth, 0, 0, 0, 0, currentDate.Location())
			if testDate.Month() == nextDueDate.Month() {
				// Date is valid for this month
				nextDueDate = testDate
				break
			} else {
				// Date rolled over to next month, so use the last day of the target month
				nextDueDate = time.Date(nextDueDate.Year(), nextDueDate.Month()+1, 1, 0, 0, 0, 0, currentDate.Location()).AddDate(0, 0, -1)
				break
			}
		}
	}

	return nextDueDate
}

// CalculateNextDueDateForExistingRecord calculates the next due date for existing records during migration
// It uses the start_date as reference and calculates the next occurrence based on current date
func CalculateNextDueDateForExistingRecord(startDate time.Time, dueDayOfMonth int, dueType string, currentDate time.Time) time.Time {
	// First, try to create a date in the current month with the specified due day
	nextDueDate := time.Date(currentDate.Year(), currentDate.Month(), dueDayOfMonth, 0, 0, 0, 0, currentDate.Location())

	// If the due day has already passed this month, move to next month
	if nextDueDate.Before(currentDate) || nextDueDate.Equal(currentDate) {
		if dueType == "monthly" {
			nextDueDate = nextDueDate.AddDate(0, 1, 0)
		} else if dueType == "yearly" {
			nextDueDate = nextDueDate.AddDate(1, 0, 0)
		} else if dueType == "weekly" {
			// For weekly, find the next occurrence
			for nextDueDate.Before(currentDate) || nextDueDate.Equal(currentDate) {
				nextDueDate = nextDueDate.AddDate(0, 0, 7)
			}
		} else if dueType == "daily" {
			// For daily, use tomorrow
			nextDueDate = currentDate.AddDate(0, 0, 1)
		}
	}

	// Handle edge case where the due day doesn't exist in the target month
	if dueType == "monthly" || dueType == "yearly" {
		for {
			// Try to create the date
			testDate := time.Date(nextDueDate.Year(), nextDueDate.Month(), dueDayOfMonth, 0, 0, 0, 0, currentDate.Location())
			if testDate.Month() == nextDueDate.Month() {
				// Date is valid for this month
				nextDueDate = testDate
				break
			} else {
				// Date rolled over to next month, so use the last day of the target month
				nextDueDate = time.Date(nextDueDate.Year(), nextDueDate.Month()+1, 1, 0, 0, 0, 0, currentDate.Location()).AddDate(0, 0, -1)
				break
			}
		}
	}

	return nextDueDate
}

// CalculateEndDate calculates the end date based on start date and due type
func CalculateEndDate(startDate time.Time, dueType string) time.Time {
	switch dueType {
	case "monthly":
		return startDate.AddDate(0, 12, 0) // 1 year from start
	case "yearly":
		return startDate.AddDate(5, 0, 0) // 5 years from start
	case "weekly":
		return startDate.AddDate(0, 3, 0) // 3 months from start
	case "daily":
		return startDate.AddDate(0, 1, 0) // 1 month from start
	default:
		return startDate.AddDate(1, 0, 0) // Default to 1 year
	}
}
