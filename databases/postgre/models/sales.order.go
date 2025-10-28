package models

import (
	"time"
)

type SalesOrder struct {
	ID            int64            `json:"id"`
	CustomerId    uint             `json:"customerId"`
	OrderDate     time.Time        `json:"orderDate"`
	Amount        float64          `json:"amount"`
	Discount      int              `json:"discount"`
	Total         float64          `json:"totalAmount"`
	Tenant        string           `json:"tenant"`
	Customer      Customers        `json:"customer" gorm:"foreignKey:ID;references:customer_id"`
	Items         []SalesOrderItem `json:"items" gorm:"Foreignkey:sales_order_id;association_foreignkey:ID;"`
}

func (c *SalesOrder) TableName() string {
	return "sales_orders"
}

type SalesOrderResponse struct {
	ID         int64               `json:"id"`
	CustomerId int64               `json:"customerId"`
	OrderDate  time.Time           `json:"orderDate"`
	Amount     float64             `json:"amount"`
	Discount   int                 `json:"discount"`
	Total      float64             `json:"totalAmount"`
	Tenant     string              `json:"tenant"`
	OrderItems []SalesOrderItemRes `json:"orderItems" gorm:"foreignKey:SalesOrderID"`
}
