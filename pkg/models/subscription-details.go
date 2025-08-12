package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Subscription_Details struct {
	bun.BaseModel         `bun:"subscription_details"`
	ID                    int       `bun:"id,pk,autoincrement" json:"id"`
	AccountID             int       `bun:"account_id" json:"account_id"`
	SubscriptionChannelID int       `bun:"subscription_channel_id" json:"subscription_channel_id"`
	StartDate             time.Time `bun:"start_date" json:"start_date"`
	DueDate               time.Time `bun:"due_date" json:"due_date"`
	Status                string    `bun:"status" json:"status"`
	StartTime             time.Time `bun:"start_time" json:"start_time"`
	DueTime               time.Time `bun:"due_time" json:"due_time"`
	MonthlyBill           float64   `bun:"monthly_bill" json:"monthly_bill"`
	ReminderDate          time.Time `bun:"reminder_date" json:"reminder_date"`
	ReminderTime          time.Time `bun:"reminder_time" json:"reminder_time"`
	CreatedAt             time.Time `bun:"created_at" json:"created_at"`
	UpdatedAt             time.Time `bun:"updated_at" json:"updated_at"`
}
