package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Subscription_Details struct {
	bun.BaseModel         `bun:"subscription_details"`
	ID                    int        `bun:"id,pk,autoincrement" json:"id"`
	AccountID             int        `bun:"account_id" json:"account_id"`
	SubscriptionChannelID int        `bun:"subscription_channel_id" json:"subscription_channel_id"`
	StartDate             time.Time  `bun:"start_date" json:"start_date"`
	NextDueDate           time.Time  `bun:"next_due_date" json:"next_due_date"`
	DueType               string     `bun:"due_type" json:"due_type"`
	DueDayOfMonth         int        `bun:"due_day_of_month" json:"due_day_of_month"`
	EndDate               *time.Time `bun:"end_date,nullzero" json:"end_date,omitempty"`
	Status                string     `bun:"status" json:"status"`
	StartTime             *time.Time `bun:"start_time,nullzero" json:"start_time,omitempty"`
	DueTime               *time.Time `bun:"due_time,nullzero" json:"due_time,omitempty"`
	MonthlyBill           float64    `bun:"monthly_bill" json:"monthly_bill"`
	ReminderDate          *time.Time `bun:"reminder_date,nullzero" json:"reminder_date,omitempty"`
	ReminderTime          *time.Time `bun:"reminder_time,nullzero" json:"reminder_time,omitempty"`
	CreatedAt             time.Time  `bun:"created_at" json:"created_at"`
	UpdatedAt             time.Time  `bun:"updated_at" json:"updated_at"`
}
