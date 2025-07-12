package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
)

// SaveProductLocation ...
func (d *postgreDatabase) SaveProductLocation(ctx context.Context, data []*models.ProductLocation) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveProductLocation")
	query := d.Db.WithContext(ctx)

	if err := query.Save(data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save Location")
		return err
	}
	return nil
}

// GetProductLocationByLocationCode ...
func (d *postgreDatabase) GetProductLocationByLocationCode(ctx context.Context,productId int64, locationCode string) (data models.ProductLocation, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("product_id = ? and location_code = ?", productId, locationCode).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetProductLocationByProductId ...
func (d *postgreDatabase) GetProductLocationByProductId(ctx context.Context,productId int64) (data []models.ProductLocation, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Order("qtty desc").Where("product_id = ?", productId).Find(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}
