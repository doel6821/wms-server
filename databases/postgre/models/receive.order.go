package models

import "time"

type ReceiveOrder struct {
	ID            int64              `json:"id"`
	InvoiceNumber string             `json:"invoiceNumber"`
	SupplierId    int64              `json:"supplierId"`
	Supplier      Supplier           `json:"supplier" gorm:"foreignKey:ID;references:supplier_id"`
	ReceiveDate   time.Time          `json:"receiveDate"`
	Status        string             `json:"status"`
	DueDate       time.Time          `json:"dueDate"`
	PaymentStatus string             `json:"paymentStatus"`
	Tenant        string             `json:"tenant"`
	Items         []ReceiveOrderItem `json:"items" gorm:"Foreignkey:receive_order_id;association_foreignkey:ID;"`
}

func (c *ReceiveOrder) TableName() string {
	return "receive_orders"
}

type ReceiveOrderResponse struct {
	ID            int64                 `json:"id"`
	InvoiceNumber string                `json:"invoiceNumber"`
	SupplierId    int64                 `json:"supplierId"`
	ReceiveDate   time.Time             `json:"receiveDate"`
	Status        string                `json:"status"`
	Tenant        string                `json:"tenant"`
	ReceiveItems  []ReceiveOrderItemRes `json:"receiveItems" gorm:"foreignKey:ReceiveOrderID"`
}
