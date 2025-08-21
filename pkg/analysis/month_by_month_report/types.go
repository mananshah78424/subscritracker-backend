package month_by_month_report

import "time"

type MonthlySubscriptionData struct {
	Month                 string    `json:"month"`
	SubscriptionChannelId int       `json:"subscription_channel_id"`
	Year                  int       `json:"year"`
	Cost                  float64   `json:"cost"`
	Status                string    `json:"status"`
	NextDueDate           time.Time `json:"next_due_date"`
}

type MonthlyReportResponse struct {
	Subscriptions []MonthlySubscriptionData `json:"subscriptions"`
	TotalCost     float64                   `json:"total_cost"`
}
