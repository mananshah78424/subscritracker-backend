package monthly_report

import (
	"log"
	"net/http"
	"subscritracker/pkg/application"
	"subscritracker/utils/account"

	"github.com/labstack/echo/v4"
)

func GetMonthlyReportHandler(c echo.Context) error {
	// Get from subscription details where account_id is the same as the account_id in the request
	app := c.Get("app").(*application.App)

	// Get user_id from JWT token and convert to int
	accountID, err := account.ConvertAccountIdStringToInt(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
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
