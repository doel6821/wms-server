package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
	"wms-server/helpers"
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

// GetSalesOrderItemByAllocation ...
func (d *postgreDatabase) GetSalesOrderItemByAllocation(ctx context.Context, ids []int) (data []models.SalesOrderItem, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("id in (" + helpers.JoinInts(ids, ",") + ") anda allocation_order_quantity > 0").Find(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


// GetSalesOrderItemBySalesOrderId ...
func (d *postgreDatabase) GetSalesOrderItemBySalesOrderId(ctx context.Context, salesOrderId, productId int64) (data models.SalesOrderItem, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("sales_order_id = ? and product_id = ?", salesOrderId, productId).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


