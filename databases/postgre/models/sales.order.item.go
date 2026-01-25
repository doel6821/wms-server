package models

type SalesOrderItem struct {
	ID                 int64     `json:"id"`
	SalesOrderId       int64     `json:"salesOrderId"`
	CustomerId         int64     `json:"customerId"`
	Customer           Customers `json:"customer" gorm:"foreignKey:ID;references:customer_id"`
	ProductId          int64     `json:"productId"`
	Product            Product   `json:"product" gorm:"foreignKey:ID;references:product_id"`
	OrderQty           int       `json:"orderQty"`
	BackOrderQty       int       `json:"backOrderQty"`
	AllocationOrderQty int       `json:"allocationOrderQty"`
	PackingOrderQty    int       `json:"packingOrderQty"`
	InvoiceOrderQty    int       `json:"invoiceOrderQty"`
	Price              float64   `json:"price"`
	Total              float64   `json:"totalAmount"`
}

func (c *SalesOrderItem) TableName() string {
	return "sales_order_items"
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

type SalesOrderItemTotal struct {
	TotalItem            int64   `json:"totalItem"`
	TotalOrderQty        int64   `json:"totalOrderQty"`
	TotalAllocationQty   int64   `json:"totalAllocationQty"`
	TotalBackOrderQty    int64   `json:"totalBackOrderQty"`
	TotalPackingOrderQty int64   `json:"totalPackingOrderQty"`
	TotalInvoiceQty      int64   `json:"totalInvoiceQty"`
	TotalAmount          float64 `json:"totalAmount"`
}
