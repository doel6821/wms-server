package models

import (
	"time"
)

type PurchaseOrder struct {
	ID            int64               `json:"id"`
	PurchaseNumber string             `json:"purchaseNumber"`
	SupplierId    int64               `json:"supplierId"`
	OrderDate     time.Time           `json:"orderDate"`
	Amount        float64             `json:"amount"`
	Discount      int                 `json:"discount"`
	Total         float64             `json:"totalAmount"`
	Tenant        string              `json:"tenant"`
	Supplier      Supplier            `json:"supplier" gorm:"foreignKey:ID;references:supplier_id"`
	Items         []PurchaseOrderItem `json:"items" gorm:"Foreignkey:purchase_order_id;association_foreignkey:ID;"`
}

func (c *PurchaseOrder) TableName() string {
	return "purchase_orders"
}

type PurchaseOrderResponse struct {
	ID         int64                  `json:"id"`
	SupplierId int64                  `json:"supplierId"`
	OrderDate  time.Time              `json:"orderDate"`
	Amount     float64                `json:"amount"`
	Discount   int                    `json:"discount"`
	Total      float64                `json:"totalAmount"`
	Tenant     string                 `json:"tenant"`
	OrderItems []PurchaseOrderItemRes `json:"orderItems" gorm:"foreignKey:PurchaseOrderID"`
}
