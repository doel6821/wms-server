package models

type RegisterLocationRequest struct {
	ID              uint64    `json:"id"`
	Code            string  `json:"code"`
	Name            string  `json:"name"`
}