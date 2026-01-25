package models

type ReqAccountReceiveble struct {
	ID              int64   `json:"id"`
	InvoiceId       int64   `json:"invoiceId"`
	Amount          float64 `json:"amount"`
	PaymentDate     string  `json:"paymentDate"`
	PaymentMethod   string  `json:"paymentMethod"`
	ReferenceNumber string  `json:"referenceNumber"`
	Notes           string  `json:"notes"`
	InvoiceNumber   string  `json:"invoiceNumber"`
}

type ReqAccountPayable struct {
	ID              int64   `json:"id"`
	ReceiveId       int64   `json:"receiveId"`
	Amount          float64 `json:"amount"`
	PaymentDate     string  `json:"paymentDate"`
	PaymentMethod   string  `json:"paymentMethod"`
	ReferenceNumber string  `json:"referenceNumber"`
	Notes           string  `json:"notes"`
	InvoiceNumber   string  `json:"invoiceNumber"`
}

type ReqListFinance struct {
	Page            int    `form:"page"`
	Limit           int    `form:"limit"`
	StartDate       string `form:"startDate"`
	EndDate         string `form:"endDate"`
	PaymentMethod   string `form:"paymentMethod"`
	ReferenceNumber string `form:"referenceNumber"`
}
