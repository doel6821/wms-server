package models

type Product struct {
	ID              uint    `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	HETPrice        float64 `json:"hetPrice"`
	CostPrice       float64 `json:"costPrice"`
	AvgPrice        float64 `json:"avgPrice"`
	StockOnHand     int     `json:"stockOnHand"`
	StockAllocation int     `json:"stockAllocation"`
	StockBackOrder  int     `json:"stockBackOrder"`
	StockPacking    int     `json:"stockPacking"`
	StockOnPurchase int     `json:"stockOnPurchase"`
	StockOnReceive  int     `json:"stockOnReceive"`
	LeadTimeDays    int     `json:"leadTimeDays"`
	Tenant          string  `json:"tenant"`
}

func (c *Product) TableName() string {
	return "product"
}
