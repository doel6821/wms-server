package models

import "time"

type PackingOrder struct {
	ID          int64     `json:"id"`
	CustomerId  int64     `json:"customerId"`
	PackingDate time.Time `json:"packingDate"`
	Status      string    `json:"status"`
	Tenant      string    `json:"tenant"`
}

func (c *PackingOrder) TableName() string {
	return "packing_order"
}

type PackingOrderResponse struct {
	ID           int64                 `json:"id"`
	CustomerId   int64                 `json:"customerId"`
	PackingDate  time.Time             `json:"packingDate"`
	Status       string                `json:"status"`
	Tenant       string                `json:"tenant"`
	PackingItems []PackingOrderItemRes `json:"packingItems" gorm:"foreignKey:PackingOrderID"`
}
