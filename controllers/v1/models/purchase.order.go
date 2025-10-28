package models

import (
	pModels "wms-server/databases/postgre/models"
)

type PurchaseOrderRequest struct {
	SupplierID    int64               `json:"supplierId"`
	SupplierName  string              `json:"supplierName"`
	Amount        float64             `json:"amount"`
	Discount      int                 `json:"discount"`
	TotalAmount   float64             `json:"totalAmount"`
	PaymentStatus string              `json:"paymentStatus"`
	OrderItems    []PurchaseOrderItem `json:"orderItems"`
}

type PurchaseOrderItem struct {
	ProductId   int     `json:"productId"`
	ProductName string  `json:"productName"`
	OrderQty    int     `json:"orderQty"`
	Price       float64 `json:"price"`
	Total       float64 `json:"totalAmount"`
}

type ReceiveOrderRequest struct {
	SupplierID    int64                `json:"supplierId"`
	InvoiceNumber string               `json:"invoiceNumber"`
	ReceiveOrders []ReceiveOrderDetail `json:"receiveOrders"`
}

type ReceiveOrderDetail struct {
	SupplierID       int64   `json:"supplierId"`
	PurchaseOrderID  int64   `json:"purchaseOrderId"`
	ProductId        int64   `json:"productId"`
	ProductName      string  `json:"productName"`
	ProductPrice     string  `json:"productPrice"`
	ReceiveOrderQtty int     `json:"receiveOrderQty"`
	PurchasePrice    float64 `json:"purchasePrice"`
}

type StockedRequest struct {
	ReceiveOrderId int64                      `json:"receiveOrderId"`
	Items          []pModels.ReceiveOrderItem `json:"items"`
}
