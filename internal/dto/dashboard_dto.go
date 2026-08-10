package dto

type DashboardResponse struct {
	TotalOrders     int64                    `json:"total_orders"`
	TotalRevenue    float64                  `json:"total_revenue"`
	PendingOrders   int64                    `json:"pending_orders"`
	CompletedOrders int64                    `json:"completed_orders"`
	MonthlyRevenue  []MonthlyRevenueResponse `json:"monthly_revenue"`
	TopProducts     []TopProductResponse     `json:"top_products"`
}

type MonthlyRevenueResponse struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
}

type TopProductResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Sold     int64  `json:"sold"`
}
