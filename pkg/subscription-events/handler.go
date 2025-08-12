package subscriptionevents

import (
	"net/http"
	"subscritracker/pkg/models"
	subscriptiondetails "subscritracker/pkg/subscription-details"
	"subscritracker/pkg/validator"

	"github.com/labstack/echo/v4"
)

func PostSubscriptionEventsHandler(c echo.Context) error {
	request, err := validator.ValidateSubscriptionEventRequest(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Validate subscription details id exists
	exists, err := subscriptiondetails.CheckSubscriptionDetailsExists(c, request.SubscriptionDetailsID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if !exists {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "subscription_details_id cannot be found"})
	}

	subscriptionEvent := models.Subscription_Event{
		SubscriptionDetailsID: request.SubscriptionDetailsID,
		AccountID:             request.AccountID,
	}

	createdSubscriptionEvent, err := CreateSubscriptionEvent(c, subscriptionEvent)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, createdSubscriptionEvent)
}
