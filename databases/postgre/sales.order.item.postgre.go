package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
)

// SaveSalesOrderItem ...
func (d *postgreDatabase) SaveSalesOrderItem(ctx context.Context, data []*models.SalesOrderItem) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveSalesOrderItem")
	query := d.Db.WithContext(ctx)

	if err := query.Save(data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save SalesOrderItem")
		return err
	}
	return nil
}


// GetList ...
func (d *postgreDatabase) GetSalesOrderItemList(ctx context.Context, SalesOrderId int64) ( []models.SalesOrderItem, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.SalesOrderItem

	err := query.Where("sales_order_id = ?", SalesOrderId).Find(&res).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get SalesOrderItem list")
		return res, err
	}

	return res, nil
}

// GetSalesOrderItemById ...
func (d *postgreDatabase) GetSalesOrderItemById(ctx context.Context, id int64) (data models.SalesOrderItem, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("id = ?", id).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


