package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
)

// SaveSupplier ...
func (d *postgreDatabase) SaveSupplier(ctx context.Context, data models.Supplier) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveSupplier")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save Supplier")
		return err
	}
	return nil
}


// GetList ...
func (d *postgreDatabase) GetSupplierList(ctx context.Context,tenant, name string, page , limit int) ( []models.Supplier, int64, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.Supplier
	var total int64

	query = query.Where("tenant = ?", tenant)
	if name != "" {
		query = query.Where("name ilike ?", "%"+name+"%")
	}

	err := query.Order("id asc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get Supplier list")
		return res, 0, err
	}

	return res, total, nil
}

// GetSupplierById ...
func (d *postgreDatabase) GetSupplierById(ctx context.Context, id int64) (data models.Supplier, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("id = ?", id).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetSupplierByPhone ...
func (d *postgreDatabase) GetSupplierByPhone(ctx context.Context, phone, tenant string) (data models.Supplier, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("Phone = ? and tenant = ?", phone, tenant).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// DeleteSupplierById ...
func (d *postgreDatabase) DeleteSupplierById(ctx context.Context, id int64)  error {
	query := d.Db.WithContext(ctx)

	data := models.Supplier{}
	
	if err := query.Where("id = ?", id).Delete(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error delete existing data")
		return err
	}
	return nil
}
