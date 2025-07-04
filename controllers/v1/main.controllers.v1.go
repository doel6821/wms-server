package controllers

import (
	v1Usecases "wms-server/usecases/v1"

	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
)

// type key string

// V1Controller ...
type (
	v1Controller struct {
		Usecase v1Usecases.Usecase
		Logs    *logs.Logger
	}
	V1Controller interface {
		HealthCheck(ctx *gin.Context)
		// BEGIN __INCLUDE_TEMPLATE__
		Login(ctx *gin.Context)
		CekToken(ctx *gin.Context)
		RegisterUser(ctx *gin.Context)
		RegisterTenant(ctx *gin.Context)

		SaveCustomer(ctx *gin.Context)
		ListCustomer(ctx *gin.Context)
		GetCustomerById(ctx *gin.Context)
		DeleteCustomerById(ctx *gin.Context)

		SaveProduct(ctx *gin.Context)
		ListProduct(ctx *gin.Context)
		GetProductById(ctx *gin.Context)
		DeleteProductById(ctx *gin.Context)
		// END __INCLUDE_TEMPLATE__
	}
)

// InitializeV1Controller ..
func InitializeV1Controller(usecases v1Usecases.Usecase, l *logs.Logger) V1Controller {
	return &v1Controller{
		Usecase: usecases,
		Logs:    l,
	}
}
