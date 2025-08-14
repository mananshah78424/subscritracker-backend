package monthly_report

import (
	"context"
	"strconv"
	"time"

	"subscritracker/pkg/application"
	"subscritracker/pkg/models"
)

func GetSubscriptionDetails(app *application.App, accountId int) ([]models.Subscription_Details, error) {
	subscriptionDetails := []models.Subscription_Details{}
	err := app.Database.NewSelect().
		Model(&subscriptionDetails).
		Where("account_id = ?", accountId).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return subscriptionDetails, nil
}

func ExtractMonthlyData(subscriptionDetails []models.Subscription_Details) ([]MonthlyData, error) {
	monthlyData := []MonthlyData{}

	for _, subscriptionDetail := range subscriptionDetails {
		month, year := ExtractMonthAndYear(subscriptionDetail.DueDate)
		summary := MonthlyData{}
		summary.Month = month
		summary.Year = year
		summary.Cost = subscriptionDetail.MonthlyBill

		monthlyData = append(monthlyData, summary)
	}

	return monthlyData, nil
}

func ExtractMonthAndYear(date time.Time) (string, string) {
	month := date.Month().String()
	year := strconv.Itoa(date.Year())
	return month, year
}

func AggregateMonthlyTotals(subscriptionDetails []models.Subscription_Details) ([]MonthlyData, error) {
	// First, get the monthly data to extract month/year info
	monthlyData, err := ExtractMonthlyData(subscriptionDetails)
	if err != nil {
		return nil, err
	}

	year := "2025"
	months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	monthlyBreakdown := []MonthlyData{}

	for _, month := range months {
		totalCost := 0.0

		// Sum up costs for this month from the monthly data
		for _, data := range monthlyData {
			if data.Month == month && data.Year == year {
				totalCost += data.Cost
			}
		}

		// Create result for this month (even if cost is 0)
		monthTotal := MonthlyData{
			Month: month,
			Year:  year,
			Cost:  totalCost,
		}

		monthlyBreakdown = append(monthlyBreakdown, monthTotal)
	}

	return monthlyBreakdown, nil
}
