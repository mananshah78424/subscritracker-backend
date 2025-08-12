package validator

import (
	"errors"
	"time"

	"github.com/labstack/echo/v4"
)

// parseTimeOfDay parses a time string in "HH:MM" format and returns a normalized time.Time
// with a consistent base date (year 0, month 1, day 1) to avoid date-related comparison issues.
//
// Benefits of this approach:
// 1. Consistent base date prevents timezone and date-related comparison issues
// 2. All time-of-day values can be compared reliably using standard time.Time methods
// 3. Database storage is consistent and predictable
// 4. Avoids issues with daylight saving time changes affecting stored times
//
// Example: "15:30" becomes 0000-01-01 15:30:00 UTC
// This allows for reliable time comparisons without worrying about the actual date.
func parseTimeOfDay(timeStr string) (*time.Time, error) {
	if timeStr == "" {
		return nil, nil
	}

	// Validate the time string format first
	if len(timeStr) != 5 || timeStr[2] != ':' {
		return nil, errors.New("time format must be HH:MM")
	}

	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return nil, errors.New("invalid time format. Expected HH:MM (e.g., 14:30)")
	}

	// Validate hour and minute ranges
	hour := t.Hour()
	minute := t.Minute()
	if hour < 0 || hour > 23 {
		return nil, errors.New("hour must be between 00 and 23")
	}
	if minute < 0 || minute > 59 {
		return nil, errors.New("minute must be between 00 and 59")
	}

	// Normalize to a fixed date (year 0, month 1, day 1) to avoid date-related issues
	normalized := time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, time.UTC)
	return &normalized, nil
}

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

	// Validate status if provided
	if request.Status != "" {
		validStatuses := []string{"active", "inactive", "paused", "cancelled"}
		statusValid := false
		for _, validStatus := range validStatuses {
			if request.Status == validStatus {
				statusValid = true
				break
			}
		}
		if !statusValid {
			return nil, errors.New("invalid status. Must be one of: active, inactive, paused, cancelled")
		}
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

	// Parse StartTime using the normalized time-of-day parser
	var err error
	if parsed.StartTime, err = parseTimeOfDay(request.StartTime); err != nil {
		return nil, errors.New("invalid start_time format. Expected HH:MM")
	}

	// Parse DueTime using the normalized time-of-day parser
	if parsed.DueTime, err = parseTimeOfDay(request.DueTime); err != nil {
		return nil, errors.New("invalid due_time format. Expected HH:MM")
	}

	// Parse ReminderDate
	if request.ReminderDate != "" {
		if t, err := time.Parse("2006-01-02", request.ReminderDate); err == nil {
			parsed.ReminderDate = &t
		} else {
			return nil, errors.New("invalid reminder_date format. Expected YYYY-MM-DD")
		}
	}

	// Parse ReminderTime using the normalized time-of-day parser
	if parsed.ReminderTime, err = parseTimeOfDay(request.ReminderTime); err != nil {
		return nil, errors.New("invalid reminder_time format. Expected HH:MM")
	}

	// Validate time logic: if both start and due times are provided, start should be before due
	if parsed.StartTime != nil && parsed.DueTime != nil {
		if parsed.StartTime.After(*parsed.DueTime) {
			return nil, errors.New("start_time cannot be after due_time")
		}
	}

	// Validate date logic: if both start and due dates are provided, start should be before due
	if parsed.StartDate != nil && parsed.DueDate != nil {
		if parsed.StartDate.After(*parsed.DueDate) {
			return nil, errors.New("start_date cannot be after due_date")
		}
	}

	return parsed, nil
}
