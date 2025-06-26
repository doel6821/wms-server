package models

import "time"

// General ...
type General struct {
	ID        *int64    `json:"id" gorm:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"`
}
