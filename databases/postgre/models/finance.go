package models

import "time"

type AccountReceiveble struct {
	ID              int64     `json:"id"`
	InvoiceId       int64     `json:"invoiceId"`
	Amount          float64   `json:"amount"`
	PaymentDate     time.Time `json:"paymentDate"`
	PaymentMethod   string    `json:"paymentMethod"`
	ReferenceNumber string    `json:"referenceNumber"`
	Notes           string    `json:"notes"`
	Tenant          string    `json:"tenant"`
	InvoiceNumber   string    `json:"invoiceNumber"`
}


func (c *AccountReceiveble) TableName() string {
	return "account_receivable"
}

type AccountPayable struct {
	ID              int64     `json:"id"`
	ReceiveId       int64     `json:"receiveId"`
	Amount          float64   `json:"amount"`
	PaymentDate     time.Time `json:"paymentDate"`
	PaymentMethod   string    `json:"paymentMethod"`
	ReferenceNumber string    `json:"referenceNumber"`
	Notes           string    `json:"notes"`
	Tenant          string    `json:"tenant"`
	InvoiceNumber   string    `json:"invoiceNumber"`
}

func (c *AccountPayable) TableName() string {
	return "account_payable"
}
