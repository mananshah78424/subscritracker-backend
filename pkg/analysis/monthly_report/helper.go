package monthly_report

import (
	"context"
	"time"

	"subscritracker/pkg/application"
	"subscritracker/pkg/models"
)

/*
**
GetSubscriptionDetails gets the subscription details for an account
**
*/
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

/*
**
ExtractMonthlyData extracts the month and year from the due date of the subscription details
and returns a list of MonthlyData objects
Creates an object like this:
[

	{
		"month": "January",
		"year": 2025,
		"cost": 100.0
	},
	{
		"month": "February",
		"year": 2025,
		"cost": 200.0
	},
	{
		"month": "January",
		"year": 2025,
		"cost": 120.0
	},
	...

]
Same month can have multiple entries if there are multiple subscriptions
**
*/
func ExtractMonthlyData(subscriptionDetails []models.Subscription_Details) ([]MonthlyData, error) {
	monthlyData := []MonthlyData{}

	for _, subscriptionDetail := range subscriptionDetails {
		month, year := ExtractMonthAndYear(subscriptionDetail.NextDueDate)
		summary := MonthlyData{}
		summary.Month = month
		summary.Year = year
		summary.Cost = subscriptionDetail.MonthlyBill

		monthlyData = append(monthlyData, summary)
	}

	return monthlyData, nil
}

/*
**
ExtractMonthAndYear extracts the month and year from the due date of the subscription details
and returns month as string and year as int
**
*/
func ExtractMonthAndYear(date time.Time) (string, int) {
	month := date.Month().String()
	year := date.Year()
	return month, year
}

/*
**
AggregateMonthlyTotals aggregates the monthly data and returns a list of MonthlyData objects
Creates an object like this:
[

	{
		"month": "January",
		"year": 2025,
		"cost": 100.0
	},
	{
		"month": "February",
		"year": 2025,
		"cost": 200.0
	},
	...

]
Each entry is the total cost for that month so total of 12 entries
**
*/
func AggregateMonthlyTotals(subscriptionDetails []models.Subscription_Details) ([]MonthlyData, error) {
	// First, get the monthly data to extract month/year info
	monthlyData, err := ExtractMonthlyData(subscriptionDetails)
	if err != nil {
		return nil, err
	}

	// Hardcoded for now
	year := 2025
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
