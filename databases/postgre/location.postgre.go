package postgre

import (
	"context"
	cModels "wms-server/controllers/v1/models"
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

// GetLocationByLocationId ...
func (d *postgreDatabase) GetLocationByLocationId(ctx context.Context,tenant string, id int64) (data models.Location, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("id = ? and tenant = ?", id, tenant).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetList ...
func (d *postgreDatabase) GetLocationList(ctx context.Context, tenant string, req cModels.LocationQueryParams) ( []models.Location, int64, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.Location
	var total int64

	if req.Code != "" {
		query = query.Where("code like ?", req.Code)
	}

	if req.Name != "" {
		query = query.Where("name like ?", req.Name)
	}

	query = query.Where("tenant = ?", tenant)

	err := query.Order("id asc").Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get Location list")
		return res, total, err
	}

	return res, total, nil
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

