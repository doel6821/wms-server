package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
)

// SaveProduct ...
func (d *postgreDatabase) SaveProduct(ctx context.Context, data models.Product) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveProduct")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save Product")
		return err
	}
	return nil
}


// GetList ...
func (d *postgreDatabase) GetProductList(ctx context.Context, tenant, name, code string, page , limit int) ( []models.Product, int64, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.Product
	var total int64

	query = query.Where("tenant = ?", tenant)
	if name != "" {
		query = query.Where("name ilike ?", "%"+name+"%")
	}

	if code != "" {
		query = query.Where("code = ?", code)
	}

	err := query.Order("id asc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(0).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get Product list")
		return res, 0, err
	}

	return res, total, nil
}

// GetProductById ...
func (d *postgreDatabase) GetProductById(ctx context.Context, id int64) (data models.Product, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("id = ?", id).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetProductByCode ...
func (d *postgreDatabase) GetProductByCode(ctx context.Context, code string) (data models.Product, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("code = ?", code).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


// DeleteProductById ...
func (d *postgreDatabase) DeleteProductById(ctx context.Context, id int64)  error {
	query := d.Db.WithContext(ctx)

	data := models.Product{}
	
	if err := query.Where("id = ?", id).Delete(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error delete existing data")
		return err
	}
	return nil
}
