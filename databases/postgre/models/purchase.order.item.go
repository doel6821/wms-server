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

type PurchaseOrderRecomendation struct {
	SupplierId  int64   `json:"supplierId"`
	ProductId   int64   `json:"productId"`
	ProductName string  `json:"productName"`
	OrderQty    int     `json:"orderQty"`
	Price       float64 `json:"price"`
	Total       float64 `json:"totalAmount"`
}

type PurchaseOrderItemTotal struct {
	TotalItem       int64   `json:"totalItem"`
	TotalOrderQty   int64   `json:"totalOrderQty"`
	TotalReceiveQty int64   `json:"totalReceiveQty"`
	TotalStockedQty int64   `json:"totalStockedQty"`
	TotalAmount     float64 `json:"totalAmount"`
}
