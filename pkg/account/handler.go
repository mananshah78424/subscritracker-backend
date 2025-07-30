package account

import (
	"log"
	"net/http"
	"strconv"
	"subscritracker/pkg/application"
	"subscritracker/pkg/models"

	"github.com/labstack/echo/v4"
)

// GetAccount returns the current user's account information
func GetAccountHandler(c echo.Context) error {
	// For now, we'll get account ID from query param
	// Later, this will come from the authenticated user's session
	accountIDStr := c.QueryParam("id")
	if accountIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Account ID is required"})
	}

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid account ID"})
	}

	app := c.Get("app").(*application.App)

	account, err := GetAccountById(app, accountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, account)
}

func GetAccountByIdHandler(c echo.Context) error {
	accountIdStr := c.Param("id")
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	app := c.Get("app").(*application.App)
	log.Println("accountId", accountId)

	account, err := GetAccountById(app, accountId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, account)
}

func UpdateAccountHandler(c echo.Context) error {
	var account models.Account

	if err := c.Bind(&account); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	app := c.Get("app").(*application.App)

	err := UpdateAccount(app, &account)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, account)
}

func GetAccountStatsHandler(c echo.Context) error {

	accountIdStr := c.QueryParam("id")
	if accountIdStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Account ID is required"})
	}

	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	app := c.Get("app").(*application.App)

	stats, err := GetAccountStats(app, accountId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, stats)
}

func CreateAccountHandler(c echo.Context) error {
	var account models.Account

	if err := c.Bind(&account); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	app := c.Get("app").(*application.App)

	err := CreateAccount(app, &account)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, account)
}
