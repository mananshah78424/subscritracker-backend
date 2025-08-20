package subscriptiondetails

import "time"

type SubscriptionDetailsByAccountID struct {
	ID                      int        `json:"id"`
	AccountID               int        `json:"account_id"`
	SubscriptionChannelID   int        `json:"subscription_channel_id"`
	SubscriptionChannelName string     `json:"subscription_channel_name"`
	ChannelImageURL         string     `json:"channel_image_url"`
	StartDate               time.Time  `json:"start_date"`
	NextDueDate             time.Time  `json:"next_due_date"`
	Status                  string     `json:"status"`
	MonthlyBill             float64    `json:"monthly_bill"`
	ReminderDate            *time.Time `bun:",nullzero" json:"reminder_date,omitempty"`
	ReminderTime            *time.Time `bun:",nullzero" json:"reminder_time,omitempty"`
}
