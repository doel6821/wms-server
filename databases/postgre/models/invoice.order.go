package models

import "time"

type InvoiceOrder struct {
	ID           uint      `json:"id"`
	SalesOrderId uint      `json:"salesOrderId"`
	InvoiceDate  time.Time `json:"invoiceDate"`
	Status       string    `json:"status"`
	Tenant       string    `json:"tenant"`
}

func (c *InvoiceOrder) TableName() string {
	return "invoice_order"
}

type InvoiceOrderResponse struct {
	ID           uint                  `json:"id"`
	SalesOrderId uint                  `json:"salesOrderId"`
	InvoiceDate  time.Time             `json:"invoiceDate"`
	Status       string                `json:"status"`
	Tenant       string                `json:"tenant"`
	InvoiceItems []InvoiceOrderItemRes `json:"invoiceItems" gorm:"foreignKey:InvoiceOrderID"`
}
