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
		
		Login(ctx *gin.Context)
		ForgotPassword(ctx *gin.Context)
		ChangePassword(ctx *gin.Context)
		CekToken(ctx *gin.Context)
		RegisterUser(ctx *gin.Context)
		UpdateUser(ctx *gin.Context)
		RegisterTenant(ctx *gin.Context)
		ListUser(ctx *gin.Context)

		GetDataDashboard(ctx *gin.Context)

		SaveCustomer(ctx *gin.Context)
		ListCustomer(ctx *gin.Context)
		GetCustomerById(ctx *gin.Context)
		DeleteCustomerById(ctx *gin.Context)

		SaveSupplier(ctx *gin.Context)
		ListSupplier(ctx *gin.Context)
		GetSupplierById(ctx *gin.Context)
		DeleteSupplierById(ctx *gin.Context)

		SaveProduct(ctx *gin.Context)
		ListProduct(ctx *gin.Context)
		GetProductById(ctx *gin.Context)
		DeleteProductById(ctx *gin.Context)
		ProductTotal(ctx *gin.Context)

		SaveLocation(ctx *gin.Context)
		ListLocation(ctx *gin.Context)
		GetLocationByCode(ctx *gin.Context)
		GetLocationById(ctx *gin.Context)
		DeleteLocationById(ctx *gin.Context)

		SaveConfiguration(ctx *gin.Context)
		ListConfiguration(ctx *gin.Context)
		GetConfigurationByName(ctx *gin.Context)
		GetConfigurationById(ctx *gin.Context)
		DeleteConfigurationById(ctx *gin.Context)

		SaveSalesOrder(ctx *gin.Context) 
		ListSalesOrder(ctx *gin.Context) 
		GetSalesOrderById(ctx *gin.Context)
		GetSalesOrderItemByProductId(ctx *gin.Context)

		SavePackingOrder(ctx *gin.Context)
		ListPackingOrder(ctx *gin.Context)
		GetPackingOrderDetailById(ctx *gin.Context)
		GetPackingOrderItemByProductId(ctx *gin.Context)

		SaveInvoice(ctx *gin.Context)
		ListInvoice(ctx *gin.Context)
		GetInvoiceDetailById(ctx *gin.Context)
		UpdateInvoiceOrder(ctx *gin.Context)

		SaveAccountPayable(ctx *gin.Context)
		ListAccountPayable(ctx *gin.Context)
		SaveAccountReceivable(ctx *gin.Context)
		ListAccountReceivable(ctx *gin.Context)

		SavePurchaseOrder(ctx *gin.Context) 
		ListPurchaseOrder(ctx *gin.Context) 
		PurchaseOrderRecomendation(ctx *gin.Context) 
		GetPurchaseOrderById(ctx *gin.Context)

		GetAvailableReceive(ctx *gin.Context)
		SaveReceiveOrder(ctx *gin.Context)
		UpdateReceiveOrder(ctx *gin.Context)
		ListReceiveOrder(ctx *gin.Context)
		GetReceiveOrderDetailById(ctx *gin.Context)
		GetReceiveOrderDetailByProductId(ctx *gin.Context)
		SaveStocked(ctx *gin.Context)
	}
)

// InitializeV1Controller ..
func InitializeV1Controller(usecases v1Usecases.Usecase, l *logs.Logger) V1Controller {
	return &v1Controller{
		Usecase: usecases,
		Logs:    l,
	}
}
