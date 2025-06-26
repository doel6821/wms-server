package models

// LogModels ...
type LogModels struct {
	METHOD   string      `json:"METHOD"`
	PATH     string      `json:"PATH"`
	HEADER   interface{} `json:"HEADER"`
	CLIENTIP string      `json:"CLIENT_IP"`
	REQUEST  interface{} `json:"REQUEST"`
	RESPONSE interface{} `json:"RESPONSE"`
	DURATION int64       `json:"DURATION"`
}
