package models

type Product struct {
	ID               int64             `json:"id"`
	Code             string            `json:"code"`
	Name             string            `json:"name"`
	HETPrice         float64           `json:"hetPrice"`
	CostPrice        float64           `json:"costPrice"`
	AvgPrice         float64           `json:"avgPrice"`
	StockOnHand      int               `json:"stockOnHand"`
	StockAllocation  int               `json:"stockAllocation"`
	StockBackOrder   int               `json:"stockBackOrder"`
	StockPacking     int               `json:"stockPacking"`
	StockOnPurchase  int               `json:"stockOnPurchase"`
	StockOnReceive   int               `json:"stockOnReceive"`
	LeadTimeDays     int               `json:"leadTimeDays"`
	Tenant           string            `json:"tenant"`
	SupplierId       int64             `json:"supplierId"`
	Supplier         Supplier          `json:"supplier" gorm:"foreignKey:ID;references:supplier_id"`
	ProductLocations []ProductLocation `json:"productLocations" gorm:"Foreignkey:product_id;association_foreignkey:ID;"`
	Demands          Demand            `json:"demands" gorm:"Foreignkey:product_id;association_foreignkey:ID;"`
}

func (c *Product) TableName() string {
	return "products"
}
