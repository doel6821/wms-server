package models

type ProductLocation struct {
	ID           int64  `json:"id"`
	ProductId    int64  `json:"productId"`
	LocationId   int64  `json:"locationId"`
	LocationCode string `json:"locationCode"`
	Qtty         int    `json:"qtty"`
}

func (c *ProductLocation) TableName() string {
	return "product_locations"
}
