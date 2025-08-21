package month_by_month_report

import (
	"net/http"
	"time"

	"subscritracker/pkg/application"
	"subscritracker/utils/account"

	"github.com/labstack/echo/v4"
)

func GetMonthByMonthHandler(c echo.Context) error {
	// Get user_id from JWT token and convert to int
	app := c.Get("app").(*application.App)

	accountID, err := account.ConvertAccountIdStringToInt(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID format"})
	}

	currentMonth := time.Now().Month()
	currentYear := time.Now().Year()

	firstDayOfMonthDateObject := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfMonthDateObject := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1)

	monthlyData, err := GetSubscriptionDetailsForMonth(app, accountID, firstDayOfMonthDateObject, lastDayOfMonthDateObject)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get subscription details for month"})
	}

	return c.JSON(http.StatusOK, monthlyData)

}
