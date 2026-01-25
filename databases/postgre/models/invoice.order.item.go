package models

type InvoiceOrderItem struct {
	ID             int64   `json:"id"`
	InvoiceOrderId int64   `json:"invoiceOrderId"`
	SalesOrderId   int64   `json:"salesOrderId"`
	ProductCode    string  `json:"productCode"`
	ProductName    string  `json:"productName"`
	Quantity       int     `json:"quantity"`
	Price          float64 `json:"price"`
	Total          float64 `json:"total"`
}

func (c *InvoiceOrderItem) TableName() string {
	return "invoice_order_items"
}

type InvoiceOrderItemRes struct {
	ID             uint    `json:"id"`
	InvoiceOrderId uint    `json:"invoiceOrderId"`
	ProductCode    string  `json:"productCode"`
	ProductName    string  `json:"productName"`
	Quantity       int     `json:"quantity"`
	Price          float64 `json:"price"`
	Total          float64 `json:"total"`
}
