package models

type RegisterLocationRequest struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type LocationQueryParams struct {
	Name  string `form:"name"`
	Code  string `form:"code"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
}
