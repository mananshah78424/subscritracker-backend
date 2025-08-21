package analysis

import (
	"subscritracker/pkg/analysis/month_by_month_report"
	"subscritracker/pkg/analysis/monthly_report"
	"subscritracker/pkg/application"
	"subscritracker/pkg/utils"
)

func RegisterRoutes(app *application.App) {
	app.Echo.GET("/v1/analysis/monthly-report", monthly_report.GetMonthlyReportHandler, utils.AuthMiddleware)
	app.Echo.GET("/v1/analysis/month-by-month-report", month_by_month_report.GetMontByMonthHandler, utils.AuthMiddleware)
}
