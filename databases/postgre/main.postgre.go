package postgre

import (
	"context"
	"fmt"
	"wms-server/databases/postgre/models"
	"wms-server/helpers"
	uModels "wms-server/usecases/v1/models"
	"time"

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
		FindTenant(ctx context.Context, tenant string) (models.User, error)
		Save(ctx context.Context, user models.User)  error

		SaveCustomers(ctx context.Context, data models.Customers) (error)
		GetCustomerList(ctx context.Context, tenant, name string, page , limit int) ( []models.Customers, int64, error)
		GetCustomerById(ctx context.Context, id int64) (data models.Customers, err error)
		GetCustomerByPhone(ctx context.Context, phone string) (data models.Customers, err error)
		DeleteCustomerById(ctx context.Context, id int64) (err error)

		SaveProduct(ctx context.Context, data models.Product) (error)
		GetProductList(ctx context.Context, tenant, name, code string, page , limit int) ( []models.Product, int64, error)
		GetProductById(ctx context.Context, id int64) (data models.Product, err error)
		GetProductByCode(ctx context.Context, code string) (data models.Product, err error)
		DeleteProductById(ctx context.Context, id int64)  error

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
