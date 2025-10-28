package models

type Customers struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Address           string `json:"address"`
	DiscountPercent   int64  `json:"discountPercent"`
	TermOfPayment     string `json:"termOfPayment"`
	CancelOnBackOrder bool   `json:"cancelOnBackOrder"`
	Tenant            string `json:"tenant"`
}

func (c *Customers) TableName() string {
	return "customers"
}
