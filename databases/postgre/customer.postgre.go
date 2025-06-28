package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
)

// SaveCustomers ...
func (d *postgreDatabase) SaveCustomers(ctx context.Context, data models.Customers) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveCustomers")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save Customers")
		return err
	}
	return nil
}


// GetList ...
func (d *postgreDatabase) GetCustomerList(ctx context.Context, name string, page , limit int) ( []models.Customers, int64, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.Customers
	var total int64

	if name != "" {
		query = query.Where("name ilike ?", "%"+name+"%")
	}

	err := query.Order("id asc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(0).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get Customers list")
		return res, 0, err
	}

	return res, total, nil
}

// GetCustomerById ...
func (d *postgreDatabase) GetCustomerById(ctx context.Context, id int64) (data models.Customers, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("id = ?", id).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetCustomerByPhone ...
func (d *postgreDatabase) GetCustomerByPhone(ctx context.Context, phone string) (data models.Customers, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("Phone = ?", phone).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// DeleteCustomerById ...
func (d *postgreDatabase) DeleteCustomerById(ctx context.Context, id int64)  error {
	query := d.Db.WithContext(ctx)

	data := models.Customers{}
	
	if err := query.Where("id = ?", id).Delete(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error delete existing data")
		return err
	}
	return nil
}
