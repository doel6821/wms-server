package helpers

import (
	"wms-server/constants"
	"wms-server/helpers/models"
	uModels "wms-server/usecases/v1/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-contrib/uuid"
)

// GenerateResponseHealthCheck ...
func GenerateResponseHealthCheck(listData ...uModels.DataHealthCheck) models.Response {
	res := models.Response{
		Meta: GetMetaResponse(constants.RC_SUCCESS),
	}
	res.Data = listData

	for _, val := range listData {
		if val.StatusCode != 200 {
			res.Meta.Code = "400"
			res.Meta.Message = "Warning"
			return res
		}
	}
	res.Meta.Message = "Healthy"
	return res
}

// GetTraceID ...
func GetTraceID(ctx *gin.Context) string {
	traceID := ctx.GetHeader("X-Trace-ID")
	if traceID == "" {
		traceID = uuid.NewV1().String()
	}
	return traceID
}

// TelnetIP ...
func TelnetIP(urlSvc string) (urlString string) {
	var urlSvc1, urlSvc2, port string
	var dataURL []string
	if strings.Contains(urlSvc, "http://") {
		urlSvc1 = strings.ReplaceAll(urlSvc, "http://", "")
		dataURL = strings.Split(urlSvc1, ":")
		if len(dataURL) > 1 {
			port = dataURL[1]
		}
		if port == "" {
			port = "80"
		}
	}

	if strings.Contains(urlSvc, "https://") {
		urlSvc2 = strings.ReplaceAll(urlSvc, "https://", "")
		dataURL = strings.Split(urlSvc2, ":")
		if len(dataURL) > 1 {
			port = dataURL[1]
		}
		if port == "" {
			port = "443"
		}
	}

	return dataURL[0] + ":" + port
}
