package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
)

// SaveLocation ...
func (d *postgreDatabase) SaveLocation(ctx context.Context, data models.Location) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveLocation")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save Location")
		return err
	}
	return nil
}



// GetLocationByLocationCode ...
func (d *postgreDatabase) GetLocationByLocationCode(ctx context.Context,tenant, locationCode string) (data models.Location, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("code = ? and tenant = ?", locationCode, tenant).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetList ...
func (d *postgreDatabase) GetLocationList(ctx context.Context, tenant string) ( []models.Location, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.Location

	err := query.Where("tenant = ?", tenant).Order("id asc").Find(&res).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get Location list")
		return res, err
	}

	return res, nil
}

// DeleteLocationById ...
func (d *postgreDatabase) DeleteLocationById(ctx context.Context, id int64)  error {
	query := d.Db.WithContext(ctx)

	data := models.Location{}
	
	if err := query.Where("id = ?", id).Delete(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error delete existing data")
		return err
	}
	return nil
}

