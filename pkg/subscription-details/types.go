package subscriptiondetails

import "time"

type SubscriptionDetailsByAccountID struct {
	ID                      int       `json:"id"`
	AccountID               int       `json:"account_id"`
	SubscriptionChannelID   int       `json:"subscription_channel_id"`
	SubscriptionChannelName string    `json:"subscription_channel_name"`
	ChannelImageURL         string    `json:"channel_image_url"`
	StartDate               time.Time `json:"start_date"`
	NextDueDate             time.Time `json:"next_due_date"`
	Status                  string    `json:"status"`
	MonthlyBill             float64   `json:"monthly_bill"`
	ReminderDate            time.Time `json:"reminder_date"`
	ReminderTime            time.Time `json:"reminder_time"`
}
