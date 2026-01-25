package postgre

import (
	"context"
	cModels "wms-server/controllers/v1/models"
	"wms-server/databases/postgre/models"
)

// SaveConfiguration ...
func (d *postgreDatabase) SaveConfiguration(ctx context.Context, data models.Configuration) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveConfiguration")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save Configuration")
		return err
	}
	return nil
}



// GetConfigurationByConfigurationName ...
func (d *postgreDatabase) GetConfigurationByConfigurationName(ctx context.Context,tenant, configurationName string) (data models.Configuration, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("name ilike '%" + configurationName +"%' and tenant = ?", tenant).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetConfigurationByConfigurationId ...
func (d *postgreDatabase) GetConfigurationByConfigurationId(ctx context.Context,tenant string, id int64) (data models.Configuration, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("id = ? and tenant = ?", id, tenant).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetList ...
func (d *postgreDatabase) GetConfigurationList(ctx context.Context, tenant string, req cModels.ConfigurationQueryParams) ( []models.Configuration, int64, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.Configuration
	var total int64

	if req.Name != "" {
		query = query.Where("name like ?", req.Name)
	}

	query = query.Where("tenant = ?", tenant)

	err := query.Order("id asc").Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get Configuration list")
		return res, total, err
	}

	return res, total, nil
}

// DeleteConfigurationById ...
func (d *postgreDatabase) DeleteConfigurationById(ctx context.Context, id int64)  error {
	query := d.Db.WithContext(ctx)

	data := models.Configuration{}
	
	if err := query.Where("id = ?", id).Delete(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error delete existing data")
		return err
	}
	return nil
}

