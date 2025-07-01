package usecases

import (
	"context"
	"wms-server/databases"
	logs "github.com/sirupsen/logrus"
	// pModels "wms-server/databases/postgre/models"
	cModels "wms-server/controllers/v1/models"
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
		RegisterUser(ctx context.Context, req cModels.RegisterRequest) hModels.Response
		RegisterTenant(ctx context.Context, req cModels.RegisterRequest) hModels.Response

		SaveCustomer(ctx context.Context, req cModels.RegisterCustomerRequest) hModels.Response
		GetCustomerList(ctx context.Context, name string, page , limit int) hModels.Response
		GetCustomerId(ctx context.Context, id int64) hModels.Response
		DeleteCustomerId(ctx context.Context, id int64) hModels.Response

		// END __INCLUDE_TEMPLATE__
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
