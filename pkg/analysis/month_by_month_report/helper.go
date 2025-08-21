package month_by_month_report

import (
	"context"
	"fmt"
	"subscritracker/pkg/application"
	"subscritracker/pkg/models"
	"time"
)

func GetSubscriptionDetailsForMonth(app *application.App, accountID int, startDate time.Time, endDate time.Time) (*MonthlyReportResponse, error) {

	var subscriptions []models.Subscription_Details
	var totalCost float64

	err := app.Database.NewSelect().
		Model(&subscriptions).
		Where("account_id = ? AND next_due_date BETWEEN ? AND ?", accountID, startDate, endDate).
		Scan(context.Background())

	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}

	var monthlySubscriptionData []MonthlySubscriptionData

	// Process each subscription
	for _, subscription := range subscriptions {
		// Create MonthlyData from the subscription
		data := MonthlySubscriptionData{
			Month:                 subscription.NextDueDate.Month().String(),
			Year:                  subscription.NextDueDate.Year(),
			Cost:                  subscription.MonthlyBill,
			SubscriptionChannelId: subscription.SubscriptionChannelID,
			Status:                subscription.Status,
			NextDueDate:           subscription.NextDueDate,
		}
		totalCost += subscription.MonthlyBill

		monthlySubscriptionData = append(monthlySubscriptionData, data)
	}

	// Create the response with both data and total cost
	response := &MonthlyReportResponse{
		Subscriptions: monthlySubscriptionData,
		TotalCost:     totalCost,
	}

	return response, nil
}
