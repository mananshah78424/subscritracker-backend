package validator

import (
	"errors"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
)

// --- Helpers & Constants ---

var (
	validStatuses   = []string{"active", "inactive", "paused", "cancelled"}
	validSortFields = map[string]string{
		"monthly_bill":  "sd.monthly_bill",
		"next_due_date": "sd.next_due_date",
		"start_date":    "sd.start_date",
		"status":        "sd.status",
		"channel_name":  "sc.channel_name",
	}
)

// parseTimeOfDay parses "HH:MM" and normalizes to 0000-01-01 HH:MM:00 UTC.
func parseTimeOfDay(timeStr string) (*time.Time, error) {
	if timeStr == "" {
		return nil, nil
	}
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return nil, errors.New("invalid time format. Expected HH:MM (e.g., 14:30)")
	}
	normalized := time.Date(0, 1, 1, t.Hour(), t.Minute(), 0, 0, time.UTC)
	return &normalized, nil
}

// parseDate parses YYYY-MM-DD format.
func parseDate(dateStr, field string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid %s format. Expected YYYY-MM-DD", field)
	}
	return &t, nil
}

// validateEnum checks if val is in allowed list.
func validateEnum(val string, allowed []string) bool {
	for _, v := range allowed {
		if val == v {
			return true
		}
	}
	return false
}

// --- Request Structs ---

type SubscriptionDetailsRequest struct {
	SubscriptionChannelID int     `json:"subscription_channel_id" form:"subscription_channel_id" validate:"required"`
	StartDate             string  `json:"start_date" form:"start_date"`
	NextDueDate           string  `json:"next_due_date" form:"next_due_date"`
	DueType               string  `json:"due_type" form:"due_type"`
	DueDayOfMonth         int     `json:"due_day_of_month" form:"due_day_of_month"`
	EndDate               string  `json:"end_date" form:"end_date"`
	Status                string  `json:"status" form:"status"`
	StartTime             string  `json:"start_time" form:"start_time"`
	DueTime               string  `json:"due_time" form:"due_time"`
	MonthlyBill           float64 `json:"monthly_bill" form:"monthly_bill"`
	ReminderDate          string  `json:"reminder_date" form:"reminder_date"`
	ReminderTime          string  `json:"reminder_time" form:"reminder_time"`
}

type ParsedSubscriptionDetails struct {
	SubscriptionChannelID int
	StartDate             *time.Time
	NextDueDate           *time.Time
	EndDate               *time.Time
	Status                string
	DueType               string
	DueDayOfMonth         int
	MonthlyBill           float64
	StartTime             *time.Time
	DueTime               *time.Time
	ReminderDate          *time.Time
	ReminderTime          *time.Time
}

type FilterOptions struct {
	Status          string     `json:"status" query:"status"`
	SortBy          string     `json:"sort_by" query:"sort_by"`
	SortOrder       string     `json:"sort_order" query:"sort_order"`
	MinCost         *float64   `json:"min_cost" query:"min_cost"`
	MaxCost         *float64   `json:"max_cost" query:"max_cost"`
	StartDateFrom   *time.Time `json:"start_date_from" query:"start_date_from"`
	StartDateTo     *time.Time `json:"start_date_to" query:"start_date_to"`
	NextDueDateFrom *time.Time `json:"next_due_date_from" query:"next_due_date_from"`
	NextDueDateTo   *time.Time `json:"next_due_date_to" query:"next_due_date_to"`
}

// --- Validators ---

func ValidateSubscriptionDetailsRequest(c echo.Context) (*ParsedSubscriptionDetails, error) {
	var req SubscriptionDetailsRequest
	if err := c.Bind(&req); err != nil {
		return nil, err
	}

	parsed := &ParsedSubscriptionDetails{
		SubscriptionChannelID: req.SubscriptionChannelID,
		MonthlyBill:           req.MonthlyBill,
		Status:                defaultIfEmpty(req.Status, "active"),
		DueType:               defaultIfEmpty(req.DueType, "monthly"),
		DueDayOfMonth:         defaultIfZero(req.DueDayOfMonth, 1),
	}

	// Validate status
	if !validateEnum(parsed.Status, validStatuses) {
		return nil, errors.New("invalid status. Must be one of: active, inactive, paused, cancelled")
	}

	// Dates
	var err error
	if parsed.StartDate, err = parseDate(req.StartDate, "start_date"); err != nil {
		return nil, err
	}
	if parsed.StartDate == nil {
		now := time.Now()
		parsed.StartDate = &now
	}
	if parsed.NextDueDate, err = parseDate(req.NextDueDate, "next_due_date"); err != nil {
		return nil, err
	}
	if parsed.EndDate, err = parseDate(req.EndDate, "end_date"); err != nil {
		return nil, err
	}
	if parsed.ReminderDate, err = parseDate(req.ReminderDate, "reminder_date"); err != nil {
		return nil, err
	}

	// Times
	if parsed.StartTime, err = parseTimeOfDay(req.StartTime); err != nil {
		return nil, err
	}
	if parsed.DueTime, err = parseTimeOfDay(req.DueTime); err != nil {
		return nil, err
	}
	if parsed.ReminderTime, err = parseTimeOfDay(req.ReminderTime); err != nil {
		return nil, err
	}

	// Channel Id
	if parsed.SubscriptionChannelID == 0 {
		return nil, errors.New("subscription_channel_id is required")
	}

	// Validate time logic
	if parsed.StartTime != nil && parsed.DueTime != nil && parsed.StartTime.After(*parsed.DueTime) {
		return nil, errors.New("start_time cannot be after due_time")
	}

	// Validate bill
	if parsed.MonthlyBill == 0 {
		return nil, errors.New("monthly_bill is required")
	}

	return parsed, nil
}

func ValidateSubscriptionDetailsFilters(c echo.Context) (*FilterOptions, error) {
	var filters FilterOptions
	if err := c.Bind(&filters); err != nil {
		return nil, errors.New("invalid filter parameters")
	}

	// Validate sort order
	if filters.SortOrder != "" && filters.SortOrder != "asc" && filters.SortOrder != "desc" {
		return nil, errors.New("sort_order must be 'asc' or 'desc'")
	}

	// Validate status
	if filters.Status != "" && !validateEnum(filters.Status, validStatuses) {
		return nil, errors.New("invalid status. Must be one of: active, inactive, paused, cancelled")
	}

	// Validate cost range
	if filters.MinCost != nil && filters.MaxCost != nil && *filters.MinCost > *filters.MaxCost {
		return nil, errors.New("min_cost cannot be greater than max_cost")
	}

	// Validate sort field
	if filters.SortBy != "" {
		if _, ok := validSortFields[filters.SortBy]; !ok {
			return nil, fmt.Errorf("invalid sort field: %s. Must be one of: monthly_bill, next_due_date, start_date, status, channel_name", filters.SortBy)
		}
	} else if filters.SortOrder != "" {
		return nil, errors.New("sort_order requires sort_by to be specified")
	}

	// Validate date ranges
	if filters.StartDateFrom != nil && filters.StartDateTo != nil && filters.StartDateFrom.After(*filters.StartDateTo) {
		return nil, errors.New("start_date_from cannot be after start_date_to")
	}
	if filters.NextDueDateFrom != nil && filters.NextDueDateTo != nil && filters.NextDueDateFrom.After(*filters.NextDueDateTo) {
		return nil, errors.New("next_due_date_from cannot be after next_due_date_to")
	}

	return &filters, nil
}

// --- Utility defaults ---

func defaultIfEmpty(val, def string) string {
	if val == "" {
		return def
	}
	return val
}
func defaultIfZero(val, def int) int {
	if val == 0 {
		return def
	}
	return val
}
