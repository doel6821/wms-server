package models

type ReceiveOrderItem struct {
	ID              int64   `json:"id"`
	ReceiveOrderId  int64   `json:"receiveOrderId"`
	PurchaseOrderId int64   `json:"purchaseOrderId"`
	ProductId       int64   `json:"productId"`
	ProductCode     string  `json:"productCode"`
	ProductName     string  `json:"productName"`
	ProductLocation string  `json:"productLocation"`
	ReceiveOrderQty int     `json:"receiveOrderQty"`
	PurchasePrice   float64 `json:"purchasePrice"`
}

func (c *ReceiveOrderItem) TableName() string {
	return "receive_order_items"
}

type ReceiveOrderItemRes struct {
	ID              int64   `json:"id"`
	ReceiveOrderId  int64   `json:"ReceiveOrderId"`
	PurchaseOrderId int64   `json:"purchaseOrderId"`
	ProductId       int64   `json:"productId"`
	ProductCode     string  `json:"productCode"`
	ProductName     string  `json:"productName"`
	ProductPrice    float64 `json:"productPrice"`
	LocationId      int64   `json:"locationId"`
	Location        string  `json:"location"`
	ReceiveOrderQty int     `json:"ReceiveOrderQty"`
}
