package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
)

// SaveDemand ...
func (d *postgreDatabase) SaveDemand(ctx context.Context, data []*models.Demand) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveDemand")
	query := d.Db.WithContext(ctx)

	if err := query.Save(data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save Demand")
		return err
	}
	return nil
}



// GetDemandByProductId ...
func (d *postgreDatabase) GetDemandByProductId(ctx context.Context, productId int64) (data models.Demand, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("product_id = ?", productId).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}
