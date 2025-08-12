package validator

import (
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

type SubscriptionEventRequest struct {
	SubscriptionDetailsID int `json:"subscription_details_id"`
	AccountID             int `json:"account_id"`
}

func ValidateSubscriptionEventRequest(c echo.Context) (SubscriptionEventRequest, error) {
	request := SubscriptionEventRequest{}

	// Try to bind JSON first
	if err := c.Bind(&request); err == nil {
		// JSON binding succeeded, validate the values
		if request.SubscriptionDetailsID <= 0 {
			return SubscriptionEventRequest{}, errors.New("subscription_details_id must be a positive integer")
		}
		if request.AccountID <= 0 {
			return SubscriptionEventRequest{}, errors.New("account_id must be a positive integer")
		}
		return request, nil
	}

	// If JSON binding failed, try form data
	subscriptionDetailsIDStr := c.FormValue("subscription_details_id")
	accountIDStr := c.FormValue("account_id")

	if subscriptionDetailsIDStr == "" {
		return SubscriptionEventRequest{}, errors.New("subscription_details_id is required")
	}

	if accountIDStr == "" {
		return SubscriptionEventRequest{}, errors.New("account_id is required")
	}

	// Parse string values to integers
	subscriptionDetailsID, err := strconv.Atoi(subscriptionDetailsIDStr)
	if err != nil {
		return SubscriptionEventRequest{}, errors.New("subscription_details_id must be a valid integer")
	}

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		return SubscriptionEventRequest{}, errors.New("account_id must be a valid integer")
	}

	request.SubscriptionDetailsID = subscriptionDetailsID
	request.AccountID = accountID

	return request, nil
}
