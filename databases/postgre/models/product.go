package models

type Product struct {
	ID              uint    `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	HETPrice        float64 `json:"het_price"`
	CostPrice       float64 `json:"cost_price"`
	AvgPrice        float64 `json:"avg_price"`
	StockOnHand     int     `json:"stock_on_hand"`
	StockAllocation int     `json:"stock_allocation"`
	StockPacking    int     `json:"stock_packing"`
	StockOnPurchase int     `json:"stock_on_purchase"`
	StockOnReceive  int     `json:"stock_on_receive"`
	LeadTimeDays    int     `json:"lead_time_days"`
}