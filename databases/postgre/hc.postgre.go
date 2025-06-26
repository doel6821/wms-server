package postgre

import (
	"context"
	"wms-server/helpers"
	uModels "wms-server/usecases/v1/models"
)

// HealthCheck database postgre ...
func (d *postgreDatabase) HealthCheck(ctx context.Context) uModels.DataHealthCheck {
	res := uModels.DataHealthCheck{
		ServiceName: "Postgres Database",
		Host:        helpers.GetEnv("HOST_POSTGRE") + ":" + helpers.GetEnv("PORT_POSTGRE"),
		StatusCode:  200,
	}
	var data string
	err := d.Db.Raw(`SELECT VERSION()::text`).Scan(&data).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting version information")
		res.StatusCode = 400
		return res
	}
	res.AdditionalData = data
	d.Logs.WithContext(ctx).Info(res)
	return res
}
