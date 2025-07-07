package models

import (
	"time"
)

type SalesOrder struct {
	ID         uint      `json:"id"`
	CustomerId uint      `json:"customerId"`
	OrderDate  time.Time `json:"orderDate"`
	Amount     float64   `json:"amount"`
	Discount   int       `json:"discount"`
	Total      float64   `json:"totalAmount"`
	Status     string    `json:"status"`
	Tenant     string    `json:"tenant"`
}

func (c *SalesOrder) TableName() string {
	return "sales_order"
}

type SalesOrderResponse struct {
	ID         uint                `json:"id"`
	CustomerId uint                `json:"customerId"`
	OrderDate  time.Time           `json:"orderDate"`
	Amount     float64             `json:"amount"`
	Discount   int                 `json:"discount"`
	Total      float64             `json:"totalAmount"`
	Status     string              `json:"status"`
	Tenant     string              `json:"tenant"`
	OrderItems []SalesOrderItemRes `json:"orderItems" gorm:"foreignKey:SalesOrderID"`
}
