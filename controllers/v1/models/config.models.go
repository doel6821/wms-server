package models

type RegisterConfigurationRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

type ConfigurationQueryParams struct {
	Name  string `form:"name"`
	Page  int    `form:"page"`
	Limit int    `form:"limit"`
}
