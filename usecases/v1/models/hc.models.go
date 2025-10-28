package models

// CheckServiceReq ..
type CheckServiceReq struct {
	ServiceName string
	Host        string
	Endpoint    string
	Method      string
}

// HealthCheckData represents health check response data
type HealthCheckData []DataHealthCheck

// ResponseMessage maintains backward compatibility but follows standard pattern
// Consider using hModels.Response with HealthCheckData in Data field instead
type ResponseMessage struct {
	ResponseCode int               `json:"rc"`
	Message      string            `json:"message"`
	Data         []DataHealthCheck `json:"data"`
}

// DataHealthCheck ..
type DataHealthCheck struct {
	ServiceName    string `json:"serviceName"`
	Host           string `json:"host"`
	StatusCode     int    `json:"statusCode"`
	AdditionalData string `json:"additionalData"`
}
