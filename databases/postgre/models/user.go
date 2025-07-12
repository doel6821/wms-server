package models

type User struct {
	ID       int64   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Tenant   string `json:"tenant"`
}

func (c *User) TableName() string {
	return "users"
}

