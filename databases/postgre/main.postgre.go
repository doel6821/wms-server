package postgre

import (
	"context"
	"fmt"
	"time"
	cModels "wms-server/controllers/v1/models"
	"wms-server/databases/postgre/models"
	"wms-server/helpers"
	uModels "wms-server/usecases/v1/models"

	logs "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgreDatabase ...
type (
	postgreDatabase struct {
		Db   *gorm.DB
		Logs *logs.Logger
	}

	PostgreDatabase interface {
		HealthCheck(ctx context.Context) uModels.DataHealthCheck
		FindUser(ctx context.Context, email string) (models.User, error)
		FindUserById(ctx context.Context, id int) (models.User, error)
		FindTenant(ctx context.Context, tenant string) (models.User, error)
		Save(ctx context.Context, user models.User) error
		GetUserList(ctx context.Context, tenant, email string, page, limit int) ([]models.User, int64, error)

		SaveCustomers(ctx context.Context, data models.Customers) error
		GetCustomerList(ctx context.Context, tenant, name string, page, limit int) ([]models.Customers, int64, error)
		GetCustomerById(ctx context.Context, id int64) (data models.Customers, err error)
		GetCustomerByPhone(ctx context.Context, phone, tenant string) (data models.Customers, err error)
		DeleteCustomerById(ctx context.Context, id int64) (err error)

		SaveSupplier(ctx context.Context, data models.Supplier) error
		GetSupplierList(ctx context.Context, tenant, name string, page, limit int) ([]models.Supplier, int64, error)
		GetSupplierById(ctx context.Context, id int64) (data models.Supplier, err error)
		GetSupplierByPhone(ctx context.Context, phone, tenant string) (data models.Supplier, err error)
		DeleteSupplierById(ctx context.Context, id int64) (err error)

		SaveProduct(ctx context.Context, data models.Product) error
		GetProductList(ctx context.Context, name, tenant, code, supplier string, page, limit int) ([]models.Product, int64, error)
		GetProductById(ctx context.Context, id int64) (data models.Product, err error)
		GetProductByCode(ctx context.Context, tenant, code string) (data models.Product, err error)
		DeleteProductById(ctx context.Context, id int64) error
		GetProductTotal(ctx context.Context, tenant string) ( models.ProductTotal,  error)

		SaveDemand(ctx context.Context, data []*models.Demand) error
		GetDemandByProductId(ctx context.Context, productId int64) (data models.Demand, err error)

		SaveLocation(ctx context.Context, data models.Location) error
		GetLocationList(ctx context.Context, tenant string, query cModels.LocationQueryParams) ([]models.Location, int64, error)
		GetLocationByLocationCode(ctx context.Context, tenant, locationCode string) (data models.Location, err error)
		GetLocationByLocationId(ctx context.Context, tenant string, id int64) (data models.Location, err error)
		DeleteLocationById(ctx context.Context, id int64) error

		SaveConfiguration(ctx context.Context, data models.Configuration) error
		GetConfigurationList(ctx context.Context, tenant string, query cModels.ConfigurationQueryParams) ([]models.Configuration, int64, error)
		GetConfigurationByConfigurationId(ctx context.Context, tenant string, id int64) (data models.Configuration, err error)
		GetConfigurationByConfigurationName(ctx context.Context,tenant, ConfigurationName string) (data models.Configuration, err error)
		DeleteConfigurationById(ctx context.Context, id int64) error

		SaveProductLocation(ctx context.Context, data []*models.ProductLocation) error
		GetProductLocationByLocationCode(ctx context.Context, productId int64, locationCode string) (data models.ProductLocation, err error)
		GetProductLocationByProductId(ctx context.Context, productId int64) (data []models.ProductLocation, err error)

		SaveSalesOrder(ctx context.Context, data models.SalesOrder) error
		GetSalesOrderList(ctx context.Context, tenant, allocation string, customerId int64, page, limit int) ([]models.SalesOrder, int64, error)
		GetSalesOrderById(ctx context.Context, id int64) (data models.SalesOrder, err error)
		GetSalesOrderItemAllocationByProductId(ctx context.Context, productId int64) (data []models.SalesOrderItem, err error)
		GetSalesOrderItemBackOrderByProductId(ctx context.Context, productId int64) (data []models.SalesOrderItem, err error)
		GetSalesOrderItemOnPackingByProductId(ctx context.Context, productId int64) (data []models.SalesOrderItem, err error)
		TxSalesOrder(ctx context.Context, reqSalesOrder models.SalesOrder, reqSalesOrderItem []models.SalesOrderItem, reqProductItems []models.Product, reqDemand []models.Demand) error

		SavePurchaseOrder(ctx context.Context, data models.PurchaseOrder) error
		GetPurchaseOrderList(ctx context.Context, tenant string, customerId int64, page, limit int) ([]models.PurchaseOrder, int64, error)
		GetPurchaseOrderById(ctx context.Context, id int64) (data models.PurchaseOrder, err error)
		TxPurchaseOrder(ctx context.Context, reqPurchaseOrder models.PurchaseOrder, reqPurchaseOrderItem []models.PurchaseOrderItem, reqProductItems []models.Product) error

		SaveSalesOrderItem(ctx context.Context, data []*models.SalesOrderItem) error
		GetSalesOrderItemList(ctx context.Context, SalesOrderId int64) ([]models.SalesOrderItem, error)
		GetSalesOrderItemById(ctx context.Context, id int64) (data models.SalesOrderItem, err error)
		GetSalesOrderItemBySalesOrderId(ctx context.Context, salesOrderId, productId int64) (data models.SalesOrderItem, err error)
		GetSalesOrderItemByAllocation(ctx context.Context, ids []int) (data []models.SalesOrderItem, err error)
		GetSalesOrderTotal(ctx context.Context, tenant string, startDate, endDate string) (data models.SalesOrderItemTotal, err error)

		SavePurchaseOrderItem(ctx context.Context, data []*models.PurchaseOrderItem) error
		GetPurchaseOrderItemList(ctx context.Context, purchaseOrderId int64) ([]models.PurchaseOrderItem, error)
		GetPurchaseOrderItemById(ctx context.Context, id int64) (data models.PurchaseOrderItem, err error)
		GetPurchaseOrderItemByPurchaseOrderId(ctx context.Context, purchaseOrderId, productId int64) (data models.PurchaseOrderItem, err error)
		GetPurchaseOrderItemByAllocation(ctx context.Context, ids []int) (data []models.PurchaseOrderItem, err error)
		GetAvailableReceiveOrder(ctx context.Context, tenant string, supplierId, productId int64) (data []models.PurchaseOrderItem, err error)
		GetPurchaseOrderTotal(ctx context.Context, tenant string, startDate, endDate string) (data models.PurchaseOrderItemTotal, err error)

		SaveReceiveOrder(ctx context.Context, data models.ReceiveOrder) error
		GetReceiveOrderItemByProductId(ctx context.Context, id int64) (data []models.ReceiveOrderItem, err error)
		GetReceiveOrderList(ctx context.Context, tenant string, customerId int64, page, limit int, dueDate, startDate, endDate, paymentStatus string) ([]models.ReceiveOrder, int64, error)
		GetReceiveOrderById(ctx context.Context, id int64) (data models.ReceiveOrder, err error)
		GetReceiveOrderByProductId(ctx context.Context, id int64) (data []models.ReceiveOrder, err error)
		UpdateStatusReceiveOrder(ctx context.Context, id int64, status string) error
		UpdatePaymentStatusReceiveOrder(ctx context.Context, id int64, status string, paymentDate time.Time) error
		TxReceiveOrder(ctx context.Context, reqReceiveOrder models.ReceiveOrder, reqReceiveOrderItem []models.ReceiveOrderItem, reqPurchaseOrderItems []models.PurchaseOrderItem, reqProducts []models.Product) error
		TxStockedOrder(ctx context.Context, reqReceiveOrder models.ReceiveOrder, reqPurchaseOrderItems []models.PurchaseOrderItem, reqProducts []models.Product, reqProductLocations []models.ProductLocation, reqReceiveOrderItems []models.ReceiveOrderItem, boSales []models.SalesOrderItem) error

		SavePackingOrder(ctx context.Context, data models.PackingOrder) error
		GetPackingOrderList(ctx context.Context, tenant string, customerId int64, page, limit int) ([]models.PackingOrder, int64, error)
		GetPackingOrderById(ctx context.Context, id int64) (data models.PackingOrder, err error)
		UpdateStatusPackingOrder(ctx context.Context, id int64, status string) error
		TxPackingOrder(ctx context.Context, reqPackingOrder models.PackingOrder, reqPackingOrderItem []models.PackingOrderItem, reqSalesOrderItems []models.SalesOrderItem, reqProducts []models.Product) error

		SavePackingOrderItem(ctx context.Context, data []*models.PackingOrderItem) error
		GetPackingOrderItemList(ctx context.Context, PackingOrderId int64) ([]models.PackingOrderItem, error)
		GetPackingOrderItemByProductId(ctx context.Context, productId int64) (data []models.PackingOrderItem, err error)
		GetPackingOrderItemById(ctx context.Context, id int64) (data models.PackingOrderItem, err error)
		UpdatePackingOrder(ctx context.Context, data models.PackingOrder) error

		SaveAccountPayable(ctx context.Context, data models.AccountPayable) error
		GetAccountPayableList(ctx context.Context, tenant string, req cModels.ReqListFinance) ([]models.AccountPayable, int64, error)
		SaveAccountReceivable(ctx context.Context, data models.AccountReceiveble) error
		GetAccountReceivableList(ctx context.Context, tenant string, req cModels.ReqListFinance) ([]models.AccountReceiveble, int64, error)
		GetAccountReceivableTotal(ctx context.Context, tenant string, req cModels.ReqListFinance) (float64, error)
		GetAccountPayableTotal(ctx context.Context, tenant string, req cModels.ReqListFinance) (float64, error)

		SaveInvoiceOrder(ctx context.Context, data models.InvoiceOrder) error
		GetInvoiceList(ctx context.Context, tenant string, customerId int64, page, limit int, dueDate, startDate, endDate, paymentStatus string) ([]models.InvoiceOrder, int64, error)
		GetInvoiceOrderById(ctx context.Context, id int64) (data models.InvoiceOrder, err error)
		UpdatePaymentStatusInvoiceOrder(ctx context.Context, id int64, status string, paymentDate time.Time) error
		TxInvoiceOrder(ctx context.Context, reqInvoiceOrder models.InvoiceOrder, reqInvoiceOrderItem []models.InvoiceOrderItem, reqPackingOrder models.PackingOrder, reqSalesOrderItem []models.SalesOrderItem, products []models.Product, productLocations []models.ProductLocation) error

		SaveInvoiceOrderItem(ctx context.Context, data []*models.InvoiceOrderItem) error
		GetInvoiceOrderItemList(ctx context.Context, InvoiceOrderId int64) ([]models.InvoiceOrderItem, error)
		GetInvoiceOrderItemById(ctx context.Context, id int64) (data models.InvoiceOrderItem, err error)
	}
)

// InitializePostgreDatabase ..
func InitializePostgreDatabase(conn *gorm.DB, log *logs.Logger) PostgreDatabase {
	return &postgreDatabase{
		Db:   conn,
		Logs: log,
	}
}

// logMode ...
var logMode = map[string]logger.LogLevel{
	"silent": logger.Silent,
	"error":  logger.Error,
	"warn":   logger.Warn,
	"info":   logger.Info,
}

// ConnectPostgre ...
func ConnectPostgre(log *logs.Logger) *gorm.DB {
	username := helpers.GetEnv("POSTGRE_USER")
	password := helpers.GetEnv("POSTGRE_PASS")
	host := helpers.GetEnv("POSTGRE_HOST")
	port := helpers.GetEnv("POSTGRE_PORT")
	dbName := helpers.GetEnv("POSTGRE_DB")
	ssl := helpers.GetEnv("POSTGRE_SSL")
	debug := helpers.GetEnv("POSTGRE_DEBUG")
	mode := helpers.GetEnv("POSTGRE_LOG_MODE")
	serviceName := helpers.GetEnv("SERVICE_NAME")

	newLogger := helpers.GormLogger(helpers.Options{
		Logger:                    log,
		LogLevel:                  logMode[mode],
		IgnoreRecordNotFoundError: false,
		SlowThreshold:             1 * time.Second,
	})

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta application_name=%s", host, username, password, dbName, port, ssl, serviceName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      newLogger,
		QueryFields: true,
	})
	if err != nil {
		log.WithContext(context.Background()).WithError(err).Error("Error open postgres connection")
		panic("Error open postgres connection")
	}

	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		logs.WithError(err).Error("Error ping postgres connection")
		panic("Error ping postgres connection")
	}

	registerCallbacks(db, log)

	logs.Info("Postgres connected successfully")

	//  Migration ===============
	// goose.SetBaseFS(embedMigrations)

	// if err := goose.SetDialect("postgres"); err != nil {
	// 	logs.WithFields(logs.Fields{"Message": err, "Goose": "postgres"}).Error(helpers.GetCaller())
	// 	return nil
	// }

	// sqlDB, _ := db.DB()

	// if err := goose.Up(sqlDB, "migrations"); err != nil {
	// 	logs.WithFields(logs.Fields{"Message": err, "Connection": dsn}).Error(helpers.GetCaller())
	// return nil
	// }
	// End Migration =============

	if debug == "true" {
		return db.Debug()
	}

	return db
}

func registerCallbacks(db *gorm.DB, log *logs.Logger) {
	// Callback untuk memeriksa koneksi sebelum setiap query
	db.Callback().Query().Before("gorm:query").Register("ping_before_query", func(tx *gorm.DB) {
		sqlDB, _ := tx.DB()
		if err := sqlDB.Ping(); err != nil {
			log.WithContext(context.Background()).Errorf("database ping error: %v", err)
			tx.Error = err
		}
	})
}
