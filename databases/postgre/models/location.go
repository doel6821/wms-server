package models

type Location struct {
	ID              int64    `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Tenant          string  `json:"tenant"`
}

func (c *Location) TableName() string {
	return "locations"
}
