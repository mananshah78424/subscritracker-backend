package monthly_report

type MonthlyData struct {
	Month string  `json:"month"`
	Year  int     `json:"year"`
	Cost  float64 `json:"cost"`
}
