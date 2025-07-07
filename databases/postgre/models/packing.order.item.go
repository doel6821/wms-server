package models

type PackingOrderItem struct {
	ID              uint   `json:"id"`
	PackingOrderId  uint   `json:"packingOrderId"`
	ProductId       uint   `json:"productId"`
	ProductCode     string `json:"productCode"`
	ProductName     string `json:"productName"`
	ProductLocation string `json:"productLocation"`
	PackingOrderQty int    `json:"packingOrderQty"`
}

func (c *PackingOrderItem) TableName() string {
	return "packing_order_item"
}

type PackingOrderItemRes struct {
	ID              uint   `json:"id"`
	PackingOrderId  uint   `json:"packingOrderId"`
	ProductId       uint   `json:"productId"`
	ProductCode     string `json:"productCode"`
	ProductName     string `json:"productName"`
	ProductLocation string `json:"productLocation"`
	PackingOrderQty int    `json:"packingOrderQty"`
}
