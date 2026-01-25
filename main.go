package main

import (
	"context"
	"os"
	"os/signal"
	"wms-server/controllers"
	v1Controllers "wms-server/controllers/v1"
	v1Usecases "wms-server/usecases/v1"
	"syscall"
	"wms-server/databases"
	"wms-server/databases/postgre"
	// "wms-server/databases/redis"
	"wms-server/helpers"
	"wms-server/routers"
	"wms-server/scheduler"
	"wms-server/host"
	
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logger := helpers.InitializeNewLogs()
	
	pConn := postgre.ConnectPostgre(logger)
	psql := postgre.InitializePostgreDatabase(pConn, logger)
	
	// rConn := redis.ConnectRedis()
	// rds := redis.InitializeRedis(rConn, logger)
	
	db := databases.InitializeDatabase(
		psql,
		// rds,
		logger,
	)

	mailSmpt:=host.InitializeMailSmpt(logger)
	
	v1Usecase := v1Usecases.InitializeV1Usecase(
		db,
		logger,
		mailSmpt,
	)
	v1Controller := v1Controllers.InitializeV1Controller(v1Usecase, logger)

	controller := controllers.InitializeController(
		v1Controller,
		logger,
	)
	router := routers.InitializeRouter(controller, logger)
	
	sch := scheduler.InitializeScheduler(v1Usecase, logger)
	sch.StartScheduler()
	
	logger.WithContext(context.Background()).Info("Finish Initializing")

	// Start Server
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- router.StartServer()
	}()

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		logger.WithContext(context.Background()).Info("got an interrupt, exiting...")
	case err := <-serverErr:
		if err != nil {
			logger.WithContext(context.Background()).Error("error while running api, exiting...", err)
		}
	}

}
