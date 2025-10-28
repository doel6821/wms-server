package models

type PackingOrderItem struct {
	ID              int64     `json:"id"`
	PackingOrderId  int64     `json:"packingOrderId"`
	SalesOrderId    int64     `json:"salesOrderId"`
	CustomerId      int64     `json:"customerId"`
	Customer        Customers `json:"customer" gorm:"foreignKey:ID;references:customer_id"`
	ProductId       int64     `json:"productId"`
	Product         Product   `json:"product" gorm:"foreignKey:ID;references:product_id"`
	ProductCode     string    `json:"productCode"`
	ProductName     string    `json:"productName"`
	ProductLocation string    `json:"productLocation"`
	PackingOrderQty int       `json:"packingOrderQty"`
}

func (c *PackingOrderItem) TableName() string {
	return "packing_order_items"
}

type PackingOrderItemRes struct {
	ID              int64   `json:"id"`
	PackingOrderId  int64   `json:"packingOrderId"`
	SalesOrderId    int64   `json:"salesOrderId"`
	ProductId       int64   `json:"productId"`
	ProductCode     string  `json:"productCode"`
	ProductName     string  `json:"productName"`
	ProductPrice    float64 `json:"productPrice"`
	LocationId      int64   `json:"locationId"`
	Location        string  `json:"location"`
	PackingOrderQty int     `json:"packingOrderQty"`
}
