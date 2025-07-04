package models

type Demand struct {
	ID        uint64 `json:"id"`
	ProductId uint64 `json:"productId"`
	Month     int64  `json:"month"`
	NQty      int64  `json:"nQty"`
	N1Qty     int64  `json:"n1Qty"`
	N2Qty     int64  `json:"n2Qty"`
	N3Qty     int64  `json:"n3Qty"`
	N4Qty     int64  `json:"n4Qty"`
	N5Qty     int64  `json:"n5Qty"`
	N6Qty     int64  `json:"n6Qty"`
	N7Qty     int64  `json:"n7Qty"`
	N8Qty     int64  `json:"n8Qty"`
	N9Qty     int64  `json:"n9Qty"`
	N10Qty    int64  `json:"n10Qty"`
	N11Qty    int64  `json:"n11Qty"`
	N12Qty    int64  `json:"n12Qty"`
}

func (c *Demand) TableName() string {
	return "demand"
}
