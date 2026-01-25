package models

type Configuration struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Tenant      string `json:"tenant"`
}

func (c *Configuration) TableName() string {
	return "configs"
}
