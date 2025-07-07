package models

type InvoiceOrderItem struct {
	ID              uint    `json:"id"`
	InvoiceOrderId  uint    `json:"invoiceOrderId"`
	ProductCode     string  `json:"productCode"`
	ProductName     string  `json:"productName"`
	ProductLocation string  `json:"productLocation"`
	Quantity        int     `json:"quantity"`
	Price           float64 `json:"price"`
	Discount        float64 `json:"discount"`
	Total           float64 `json:"total"`
}

func (c *InvoiceOrderItem) TableName() string {
	return "invoice_order_item"
}

type InvoiceOrderItemRes struct {
	ID              uint    `json:"id"`
	InvoiceOrderId  uint    `json:"invoiceOrderId"`
	ProductCode     string  `json:"productCode"`
	ProductName     string  `json:"productName"`
	ProductLocation string  `json:"productLocation"`
	Quantity        int     `json:"quantity"`
	Price           float64 `json:"price"`
	Discount        float64 `json:"discount"`
	Total           float64 `json:"total"`
}
