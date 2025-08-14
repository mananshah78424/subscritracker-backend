package analysis

import (
	"subscritracker/pkg/analysis/monthly_report"
	"subscritracker/pkg/application"
	"subscritracker/pkg/utils"
)

func RegisterRoutes(app *application.App) {
	app.Echo.GET("/v1/analysis/monthly-report", monthly_report.GetMonthlyReportHandler, utils.AuthMiddleware)
}
