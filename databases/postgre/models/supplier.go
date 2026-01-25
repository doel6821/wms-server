package models

type Supplier struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Address         string `json:"address"`
	DiscountPercent int64  `json:"discountPercent"`
	LeadTimeDays    int    `json:"leadTimeDays"`
	TermOfPayment   string `json:"termOfPayment"`
	Tenant          string `json:"tenant"`
}

func (c *Supplier) TableName() string {
	return "suppliers"
}
