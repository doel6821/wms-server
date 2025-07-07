package models

import "time"

type PackingOrder struct {
	ID           uint      `json:"id"`
	SalesOrderId uint      `json:"salesOrderId"`
	PackingDate  time.Time `json:"packingDate"`
	Status       string    `json:"status"`
	Tenant       string    `json:"tenant"`
}

func (c *PackingOrder) TableName() string {
	return "packing_order"
}

type PackingOrderResponse struct {
	ID           uint                  `json:"id"`
	SalesOrderId uint                  `json:"salesOrderId"`
	PackingDate  time.Time             `json:"packingDate"`
	Status       string                `json:"status"`
	Tenant       string                `json:"tenant"`
	PackingItems []PackingOrderItemRes `json:"packingItems" gorm:"foreignKey:PackingOrderID"`
}
