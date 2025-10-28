package models

import "time"

type PackingOrder struct {
	ID          int64              `json:"id"`
	CustomerId  int64              `json:"customerId"`
	Customer    Customers          `json:"customer" gorm:"foreignKey:ID;references:customer_id"`
	PackingDate time.Time          `json:"packingDate"`
	Status      string             `json:"status"`
	Tenant      string             `json:"tenant"`
	Items       []PackingOrderItem `json:"items" gorm:"Foreignkey:packing_order_id;association_foreignkey:ID;"`
}

func (c *PackingOrder) TableName() string {
	return "packing_orders"
}

type PackingOrderResponse struct {
	ID           int64                 `json:"id"`
	CustomerId   int64                 `json:"customerId"`
	PackingDate  time.Time             `json:"packingDate"`
	Status       string                `json:"status"`
	Tenant       string                `json:"tenant"`
	PackingItems []PackingOrderItemRes `json:"packingItems" gorm:"foreignKey:PackingOrderID"`
}
