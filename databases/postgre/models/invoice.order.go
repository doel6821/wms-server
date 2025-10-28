package models

import "time"

type InvoiceOrder struct {
	ID             int64              `json:"id"`
	InvoiceNumber  string             `json:"invoiceNumber"`
	PackingOrderId int64              `json:"packingOrderId"`
	CustomerId     int64              `json:"customerId"`
	CustomerName   string             `json:"customerName"`
	Customer       Customers          `json:"customer" gorm:"foreignKey:ID;references:customer_id"`
	InvoiceDate    time.Time          `json:"invoiceDate"`
	Amount         float64            `json:"amount"`
	Discount       int                `json:"discount"`
	TotalAmount    float64            `json:"totalAmount"`
	Status         string             `json:"status"`
	DueDate        time.Time          `json:"dueDate"`
	PaymentStatus  string             `json:"paymentStatus"`
	Tenant         string             `json:"tenant"`
	InvoiceItems   []InvoiceOrderItem `json:"invoiceItems" gorm:"Foreignkey:invoice_order_id;association_foreignkey:ID;"`
}

func (c *InvoiceOrder) TableName() string {
	return "invoice_orders"
}

type InvoiceOrderResponse struct {
	ID             int64                 `json:"id"`
	PackingOrderId int64                 `json:"packingOrderId"`
	CustomerId     int64                 `json:"customerId"`
	CustomerName   string                `json:"customerName"`
	InvoiceDate    time.Time             `json:"invoiceDate"`
	Status         string                `json:"status"`
	Tenant         string                `json:"tenant"`
	InvoiceItems   []InvoiceOrderItemRes `json:"invoiceItems" gorm:"foreignKey:InvoiceOrderID"`
}
