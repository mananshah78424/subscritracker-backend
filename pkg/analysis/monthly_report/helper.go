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
ExtractMonthlyData calculates monthly expenses based on subscription start date and due day
and returns a list of MonthlyData objects for the specified year
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
		"cost": 100.0
	},
	...

]
Each month gets an entry if the subscription is active during that month
**
*/
func ExtractMonthlyData(subscriptionDetails []models.Subscription_Details, targetYear int) ([]MonthlyData, error) {
	monthlyData := []MonthlyData{}

	for _, subscriptionDetail := range subscriptionDetails {
		// Calculate which months this subscription is active
		activeMonths := CalculateActiveMonths(subscriptionDetail, targetYear)

		// Add an entry for each active month
		for _, monthInfo := range activeMonths {
			summary := MonthlyData{
				Month: monthInfo.Month,
				Year:  monthInfo.Year,
				Cost:  subscriptionDetail.MonthlyBill,
			}
			monthlyData = append(monthlyData, summary)
		}
	}

	return monthlyData, nil
}

/*
**
CalculateActiveMonths determines which months a subscription is active based on:
- Start date: When subscription began
- Due day: Day of month when payment is due
- Target year: Which year to calculate for
**
*/
func CalculateActiveMonths(subscription models.Subscription_Details, targetYear int) []MonthlyData {
	var activeMonths []MonthlyData

	// Get start date components
	startYear := subscription.StartDate.Year()
	startMonth := subscription.StartDate.Month()
	startDay := subscription.StartDate.Day()

	// If subscription hasn't started yet, return empty
	if startYear > targetYear {
		return activeMonths
	}

	// Determine which months to include
	months := []time.Month{
		time.January, time.February, time.March, time.April, time.May, time.June,
		time.July, time.August, time.September, time.October, time.November, time.December,
	}

	for _, month := range months {
		// Check if subscription is active in this month
		if IsSubscriptionActiveInMonth(subscription, targetYear, month, startYear, startMonth, startDay) {
			monthName := month.String()
			activeMonths = append(activeMonths, MonthlyData{
				Month: monthName,
				Year:  targetYear,
				Cost:  subscription.MonthlyBill,
			})
		}
	}

	return activeMonths
}

/*
**
IsSubscriptionActiveInMonth checks if a subscription is active in a specific month
**
*/
func IsSubscriptionActiveInMonth(subscription models.Subscription_Details, targetYear int, targetMonth time.Month, startYear int, startMonth time.Month, startDay int) bool {
	// If subscription hasn't started yet, it's not active
	if targetYear < startYear || (targetYear == startYear && targetMonth < startMonth) {
		return false
	}

	// If subscription started in this month, check if it started before the due day
	if targetYear == startYear && targetMonth == startMonth {
		// Only active if start day is before or equal to due day
		return startDay <= subscription.DueDay
	}

	// For all other months, subscription is active if it started before this month
	return true
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
	monthlyData, err := ExtractMonthlyData(subscriptionDetails, 2025) // Assuming targetYear is 2025 for now
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
