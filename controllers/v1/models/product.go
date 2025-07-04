package models

type RegisterProductRequest struct {
	ID              uint    `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	HETPrice        float64 `json:"hetPrice"`
	CostPrice       float64 `json:"costPrice"`
	AvgPrice        float64 `json:"avgPrice"`
	StockOnHand     int     `json:"stockOnHand"`
	StockAllocation int     `json:"stockAllocation"`
	StockPacking    int     `json:"stockPacking"`
	StockOnPurchase int     `json:"stockOnPurchase"`
	StockOnReceive  int     `json:"stockOnReceive"`
	LeadTimeDays    int     `json:"leadTimeDays"`
	Tenant          string `json:"tenant"`
}

