package usecases

import (
	"context"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
)

// HealthCheck ...
func (u *usecase) HealthCheck(ctx context.Context) hModels.Response {
	// rHc := u.DB.GetRedis().HeatchCheck(ctx)
	dHc := u.DB.GetPostgre().HealthCheck(ctx)
	
	res := helpers.GenerateResponseHealthCheck(
		dHc,
		// rHc,
	)
	return res
}
