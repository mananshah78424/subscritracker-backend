package validator

import (
	"errors"
	"time"

	"github.com/labstack/echo/v4"
)

type SubscriptionDetailsRequest struct {
	SubscriptionChannelID int     `json:"subscription_channel_id" form:"subscription_channel_id" validate:"required"`
	StartDate             string  `json:"start_date" form:"start_date"`
	DueDate               string  `json:"due_date" form:"due_date"`
	Status                string  `json:"status" form:"status"`
	StartTime             string  `json:"start_time" form:"start_time"`
	DueTime               string  `json:"due_time" form:"due_time"`
	MonthlyBill           float64 `json:"monthly_bill" form:"monthly_bill"`
	ReminderDate          string  `json:"reminder_date" form:"reminder_date"`
	ReminderTime          string  `json:"reminder_time" form:"reminder_time"`
}

// ParsedSubscriptionDetails contains the parsed time.Time values
type ParsedSubscriptionDetails struct {
	SubscriptionChannelID int
	StartDate             *time.Time
	DueDate               *time.Time
	Status                string
	StartTime             *time.Time
	DueTime               *time.Time
	MonthlyBill           float64
	ReminderDate          *time.Time
	ReminderTime          *time.Time
}

func ValidateSubscriptionDetailsRequest(c echo.Context) (*ParsedSubscriptionDetails, error) {
	var request SubscriptionDetailsRequest
	if err := c.Bind(&request); err != nil {
		return nil, err
	}

	parsed := &ParsedSubscriptionDetails{
		SubscriptionChannelID: request.SubscriptionChannelID,
		Status:                request.Status,
		MonthlyBill:           request.MonthlyBill,
	}

	// Parse StartDate
	if request.StartDate != "" {
		if t, err := time.Parse("2006-01-02", request.StartDate); err == nil {
			parsed.StartDate = &t
		} else {
			return nil, errors.New("invalid start_date format. Expected YYYY-MM-DD")
		}
	}

	// Parse DueDate
	if request.DueDate != "" {
		if t, err := time.Parse("2006-01-02", request.DueDate); err == nil {
			parsed.DueDate = &t
		} else {
			return nil, errors.New("invalid due_date format. Expected YYYY-MM-DD")
		}
	}

	// Parse StartTime
	if request.StartTime != "" {
		if t, err := time.Parse("15:04", request.StartTime); err == nil {
			parsed.StartTime = &t
		} else {
			return nil, errors.New("invalid start_time format. Expected HH:MM")
		}
	}

	// Parse DueTime
	if request.DueTime != "" {
		if t, err := time.Parse("15:04", request.DueTime); err == nil {
			parsed.DueTime = &t
		} else {
			return nil, errors.New("invalid due_time format. Expected HH:MM")
		}
	}

	// Parse ReminderDate
	if request.ReminderDate != "" {
		if t, err := time.Parse("2006-01-02", request.ReminderDate); err == nil {
			parsed.ReminderDate = &t
		} else {
			return nil, errors.New("invalid reminder_date format. Expected YYYY-MM-DD")
		}
	}

	// Parse ReminderTime
	if request.ReminderTime != "" {
		if t, err := time.Parse("15:04", request.ReminderTime); err == nil {
			parsed.ReminderTime = &t
		} else {
			return nil, errors.New("invalid reminder_time format. Expected HH:MM")
		}
	}

	return parsed, nil
}
