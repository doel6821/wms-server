package models

type SalesOrderItem struct {
	ID                 uint    `json:"id"`
	SalesOrderId       uint    `json:"salesOrderId"`
	ProductId          uint    `json:"productId"`
	OrderQty           int     `json:"orderQty"`
	BackOrderQty       int     `json:"backOrderQty"`
	AllocationOrderQty int     `json:"allocationOrderQty"`
	PackingOrderQty    int     `json:"packingOrderQty"`
	InvoiceOrderQty    int     `json:"invoiceOrderQty"`
	Price              float64 `json:"price"`
	Total              float64 `json:"totalAmount"`
}

func (c *SalesOrderItem) TableName() string {
	return "sales_order_item"
}

type SalesOrderItemRes struct {
	ID                 uint    `json:"id"`
	SalesOrderId       uint    `json:"salesOrderId"`
	ProductId          uint    `json:"productId"`
	Product            Product `json:"product" gorm:"foreignKey:ProductID"`
	OrderQty           int     `json:"orderQty"`
	BackOrderQty       int     `json:"backOrderQty"`
	AllocationOrderQty int     `json:"allocationOrderQty"`
	PackingOrderQty    int     `json:"packingOrderQty"`
	InvoiceOrderQty    int     `json:"invoiceOrderQty"`
	Price              float64 `json:"price"`
	Total              float64 `json:"totalAmount"`
}
