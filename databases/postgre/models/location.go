package models

type Location struct {
	ID              uint64    `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
	Tenant          string  `json:"tenant"`
}

func (c *Location) TableName() string {
	return "location"
}
