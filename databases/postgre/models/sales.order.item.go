package models

type SalesOrderItem struct {
	ID                 int64   `json:"id"`
	SalesOrderId       int64   `json:"salesOrderId"`
	ProductId          int64   `json:"productId"`
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
	ID                 int64   `json:"id"`
	SalesOrderId       int64   `json:"salesOrderId"`
	ProductId          int64   `json:"productId"`
	Product            Product `json:"product" gorm:"foreignKey:ProductID"`
	OrderQty           int     `json:"orderQty"`
	BackOrderQty       int     `json:"backOrderQty"`
	AllocationOrderQty int     `json:"allocationOrderQty"`
	PackingOrderQty    int     `json:"packingOrderQty"`
	InvoiceOrderQty    int     `json:"invoiceOrderQty"`
	Price              float64 `json:"price"`
	Total              float64 `json:"totalAmount"`
}
