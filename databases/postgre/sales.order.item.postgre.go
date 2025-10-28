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
	
	if err = query.Where("id in (" + helpers.JoinInts(ids, ",") + ") and allocation_order_qty > 0").Find(&data).Error; err != nil {
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


// GetSalesOrderItemAllocationByProductId ...
func (d *postgreDatabase) GetSalesOrderItemAllocationByProductId(ctx context.Context, productId int64) (data []models.SalesOrderItem, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("product_id = ? and allocation_order_qty > 0", productId).Preload("Customer").Preload("Product").Find(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetSalesOrderItemBackOrderByProductId ...
func (d *postgreDatabase) GetSalesOrderItemBackOrderByProductId(ctx context.Context, productId int64) (data []models.SalesOrderItem, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Table("sales_order_items soi").Joins("left join sales_orders so on so.id = soi.sales_order_id").Where("product_id = ? and back_order_qty > 0", productId).Preload("Customer").Preload("Product").Order("so.order_date asc").Find(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


// GetSalesOrderItemOnPackingByProductId ...
func (d *postgreDatabase) GetSalesOrderItemOnPackingByProductId(ctx context.Context, productId int64) (data []models.SalesOrderItem, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Table("sales_order_items soi").Joins("left join sales_orders so on so.id = soi.sales_order_id").Where("product_id = ? and packing_order_qty > 0", productId).Preload("Customer").Preload("Product").Order("so.order_date ASC").Find(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}



