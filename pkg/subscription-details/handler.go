package subscriptiondetails

import (
	"log"
	"net/http"
	"strconv"
	"subscritracker/pkg/application"
	"subscritracker/pkg/models"
	subscription_channels "subscritracker/pkg/subscription-channels"
	"subscritracker/pkg/validator"
	"time"

	"github.com/labstack/echo/v4"
)

func PostSubscriptionDetailsHandler(c echo.Context) error {
	accountID := c.Get("user_id").(int)
	request, err := validator.ValidateSubscriptionDetailsRequest(c)
	if err != nil {
		log.Println("Error validating subscription details request:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Check if subscription channel exists
	app := c.Get("app").(*application.App)
	subscriptionChannel, err := subscription_channels.GetChannelById(c, strconv.Itoa(request.SubscriptionChannelID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get subscription channel: " + err.Error()})
	}

	// Check if user already has a subscription to this channel
	existingSubscription, err := CheckExistingSubscriptionByChannel(app, accountID, subscriptionChannel.ID)
	if err != nil {
		log.Println("Error checking existing subscription:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check existing subscription"})
	}

	if existingSubscription {
		return c.JSON(http.StatusConflict, map[string]string{"error": "You already have a subscription to this channel"})
	}

	subscriptionDetails := models.Subscription_Details{
		AccountID:             accountID,
		SubscriptionChannelID: request.SubscriptionChannelID,
		Status:                request.Status,
		MonthlyBill:           request.MonthlyBill,
		DueType:               request.DueType,
		DueDayOfMonth:         request.DueDayOfMonth,
	}

	// Handle optional time fields - only set if they exist
	if request.StartDate != nil {
		subscriptionDetails.StartDate = *request.StartDate
	}
	if request.NextDueDate != nil {
		subscriptionDetails.NextDueDate = *request.NextDueDate
	}
	if request.EndDate != nil {
		subscriptionDetails.EndDate = *request.EndDate
	}
	if request.StartTime != nil {
		subscriptionDetails.StartTime = *request.StartTime
	}
	if request.DueTime != nil {
		subscriptionDetails.DueTime = *request.DueTime
	}
	if request.ReminderDate != nil {
		subscriptionDetails.ReminderDate = *request.ReminderDate
	}
	if request.ReminderTime != nil {
		subscriptionDetails.ReminderTime = *request.ReminderTime
	}

	// Calculate NextDueDate only if it's not provided in the request
	if request.NextDueDate == nil {
		subscriptionDetails.NextDueDate = CalculateNextDueDate(subscriptionDetails.DueType, subscriptionDetails.DueDayOfMonth, time.Now())
	}

	createdSubscriptionDetails, err := CreateSubscriptionDetails(c, subscriptionDetails)
	if err != nil {
		log.Println("Error creating subscription details:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, createdSubscriptionDetails)
}

func GetSubscriptionDetailsHandler(c echo.Context) error {
	return nil
}

func GetUserSubscriptionDetailsHandler(c echo.Context) error {
	app := c.Get("app").(*application.App)
	accountID := c.Get("user_id").(int)

	// Parse and validate filter options from query parameters
	filters, err := validator.ValidateSubscriptionDetailsFilters(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	subscriptionDetails, err := GetSubscriptionDetailsByUserIdWithFilters(app, accountID, filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get subscription details"})
	}

	return c.JSON(http.StatusOK, subscriptionDetails)
}
