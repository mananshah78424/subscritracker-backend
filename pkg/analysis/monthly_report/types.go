package monthly_report

type MonthlyData struct {
	Month string  `json:"month"`
	Year  string  `json:"year"`
	Cost  float64 `json:"cost"`
}
