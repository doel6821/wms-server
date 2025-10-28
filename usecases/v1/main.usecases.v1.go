package usecases

import (
	"context"
	"wms-server/databases"
	logs "github.com/sirupsen/logrus"
	// pModels "wms-server/databases/postgre/models"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"
	hModels "wms-server/helpers/models"
	
)

// Usecase ...
type (
	usecase struct {
		DB databases.Database
		Logs *logs.Logger
	}

	Usecase interface {
		HealthCheck(ctx context.Context) hModels.Response
		Login(ctx context.Context, req cModels.LoginRequest) hModels.Response
		RegisterUser(ctx context.Context, req cModels.RegisterRequest, tenant string) hModels.Response
		UpdateUser(ctx context.Context, req cModels.RegisterRequest, tenant string, id int) hModels.Response
		RegisterTenant(ctx context.Context, req cModels.RegisterRequest) hModels.Response
		GetUserList(ctx context.Context, tenant, name string, page , limit int) hModels.Response

		SaveCustomer(ctx context.Context, req cModels.RegisterCustomerRequest) hModels.Response
		GetCustomerList(ctx context.Context, tenant, name string, page , limit int) hModels.Response
		GetCustomerId(ctx context.Context, id int64) hModels.Response
		DeleteCustomerId(ctx context.Context, id int64) hModels.Response

		SaveSupplier(ctx context.Context, req cModels.RegisterSupplierRequest) hModels.Response
		GetSupplierList(ctx context.Context, tenant, name string, page , limit int) hModels.Response
		GetSupplierId(ctx context.Context, id int64) hModels.Response
		DeleteSupplierId(ctx context.Context, id int64) hModels.Response

		SaveProduct(ctx context.Context, tenant string , req cModels.RegisterProductRequest) hModels.Response
		GetProductList(ctx context.Context,tenant, name, code, supplier string, page, limit int) hModels.Response
		GetProductId(ctx context.Context, id int64) hModels.Response
		DeleteProductId(ctx context.Context, id int64) hModels.Response 

		SaveLocation(ctx context.Context, tenant string, req cModels.RegisterLocationRequest) hModels.Response
		GetLocationList(ctx context.Context, tenant string, query cModels.LocationQueryParams) hModels.Response
		GetLocationByCode(ctx context.Context, tenant, code string) hModels.Response
		GetLocationById(ctx context.Context, tenant string, id int64) hModels.Response
		DeleteLocationId(ctx context.Context, id int64) hModels.Response

		CreateSalesOrder(ctx context.Context, tenant string, req cModels.SalesOrderRequest) hModels.Response
		GetSalesOrderList(ctx context.Context,tenant, allocation string, customerId int64, page, limit int) hModels.Response
		GetSalesOrderDetailById(ctx context.Context, id int64) hModels.Response
		GetSalesOrderItemByProductId(ctx context.Context, productId int64, tipe string) hModels.Response
		GetSalesOrderAllocationByProductId(ctx context.Context, productId int64) hModels.Response
		GetSalesOrderBackOrderByProductId(ctx context.Context, productId int64) hModels.Response 


		CreatePackingOrder(ctx context.Context, tenant string, req cModels.PackingOrderRequest) hModels.Response
		GetPackingOrderList(ctx context.Context,tenant string, customerId int64, page, limit int) hModels.Response
		GetPackingOrderDetailById(ctx context.Context, id int64) hModels.Response
		GetPackingOrderListByProductId(ctx context.Context, id int64) hModels.Response

		CreateInvoice(ctx context.Context, tenant string, req cModels.InvoiceRequest) hModels.Response 
		UpdateInvoice(ctx context.Context, tenant string, req pModels.InvoiceOrder) hModels.Response 
		GetInvoiceList(ctx context.Context,tenant string, customerId int64, page, limit int, dueDate string) hModels.Response
		GetInvoiceDetailById(ctx context.Context, id int64) hModels.Response

		CreatePurchaseOrder(ctx context.Context, tenant string, req cModels.PurchaseOrderRequest) hModels.Response
		GetPurchaseOrderList(ctx context.Context,tenant string, customerId int64, page, limit int) hModels.Response
		GetPurchaseOrderDetailById(ctx context.Context, id int64) hModels.Response

		GetAvailableReceiveList(ctx context.Context, tenant string, supplierId,productId int64) hModels.Response
		CreateReceive(ctx context.Context, tenant string, req cModels.ReceiveOrderRequest) hModels.Response
		UpdateReceive(ctx context.Context, tenant string, req pModels.ReceiveOrder) hModels.Response
		GetReceiveList(ctx context.Context,tenant string, customerId int64, page, limit int, dueDate string) hModels.Response
		GetReceiveDetailById(ctx context.Context, id int64) hModels.Response
		GetReceiveDetailByProductId(ctx context.Context, id int64) hModels.Response
		CreateStocked(ctx context.Context, tenant string, req cModels.StockedRequest) hModels.Response
	}
)

// InitializeV1Usecase ...
func InitializeV1Usecase(
	db databases.Database,
	l *logs.Logger,
) Usecase {
	uc := &usecase{
		DB:   db,
		Logs: l,
	}

	return uc
}
