package models

type SalesOrderRequest struct {
	CustomerID   int64            `json:"customerId"`
	CustomerName string           `json:"customerName"`
	Amount       float64          `json:"amount"`
	Discount     int              `json:"discount"`
	TotalAmount  float64          `json:"totalAmount"`
	OrderItems   []SalesOrderItem `json:"orderItems"`
}

type SalesOrderItem struct {
	ProductId int     `json:"productId"`
	OrderQty  int     `json:"orderQty"`
	Price     float64 `json:"price"`
	Total     float64 `json:"totalAmount"`
}

type PackingOrderRequest struct {
	CustomerID    int64   `json:"customerId"`
	SalesOrderIDs []int `json:"salesOrderIds"`
}

type InvoiceRequest struct {
	PackingOrderId int64 `json:"packingOrderId"`
}
