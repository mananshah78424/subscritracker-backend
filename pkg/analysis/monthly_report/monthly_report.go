package monthly_report

import (
	"log"
	"net/http"
	"strconv"
	"subscritracker/pkg/application"

	"github.com/labstack/echo/v4"
)

func GetMonthlyReportHandler(c echo.Context) error {
	// Get from subscription details where account_id is the same as the account_id in the request
	app := c.Get("app").(*application.App)

	// Get user_id from JWT token and convert to int
	userIDInterface := c.Get("user_id")
	var accountID int

	switch v := userIDInterface.(type) {
	case int:
		accountID = v
	case float64:
		accountID = int(v)
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			accountID = parsed
		} else {
			log.Printf("Failed to parse user_id from string: %v", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
		}
	default:
		log.Printf("Unexpected user_id type: %T, value: %v", userIDInterface, userIDInterface)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID type"})
	}

	subscriptionDetails, err := GetSubscriptionDetails(app, accountID)
	if err != nil {
		log.Printf("Error getting subscription details: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get subscription details"})
	}

	monthlyBreakdown, err := AggregateMonthlyTotals(subscriptionDetails)
	if err != nil {
		log.Printf("Error aggregating monthly totals: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to aggregate monthly totals"})
	}
	return c.JSON(http.StatusOK, monthlyBreakdown)
}
