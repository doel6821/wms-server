package models

type PurchaseOrderItem struct {
	ID              int64    `json:"id"`
	SupplierId      int64    `json:"supplierId"`
	Supplier        Supplier `json:"supplier" gorm:"foreignKey:ID;references:supplier_id"`
	PurchaseOrderId int64    `json:"purchaseOrderId"`
	ProductId       int64    `json:"productId"`
	Product         Product  `json:"product" gorm:"foreignKey:ID;references:product_id"`
	OrderQty        int      `json:"orderQty"`
	ReceiveOrderQty int      `json:"receiveOrderQty"`
	StockedOrderQty int      `json:"stockedOrderQty"`
	Price           float64  `json:"price"`
	SubTotal        float64  `json:"subTotal"`
	Discount        float64  `json:"discount"`
	Total           float64  `json:"totalAmount"`
}

func (c *PurchaseOrderItem) TableName() string {
	return "purchase_order_items"
}

type PurchaseOrderItemRes struct {
	ID              int64   `json:"id"`
	SupplierId      int64   `json:"supplierId"`
	PurchaseOrderId int64   `json:"purchaseOrderId"`
	ProductId       int64   `json:"productId"`
	Product         Product `json:"product" gorm:"foreignKey:ProductID"`
	OrderQty        int     `json:"orderQty"`
	ReceiveOrderQty int     `json:"receiveOrderQty"`
	StockedOrderQty int     `json:"stockedOrderQty"`
	Price           float64 `json:"price"`
	Total           float64 `json:"totalAmount"`
}
