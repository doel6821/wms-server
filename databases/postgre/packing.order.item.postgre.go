package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
)

// SavePackingOrderItem ...
func (d *postgreDatabase) SavePackingOrderItem(ctx context.Context, data []*models.PackingOrderItem) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SavePackingOrderItem")
	query := d.Db.WithContext(ctx)

	if err := query.Save(data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save PackingOrderItem")
		return err
	}
	return nil
}


// GetList ...
func (d *postgreDatabase) GetPackingOrderItemList(ctx context.Context, PackingOrderId int64) ( []models.PackingOrderItem, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.PackingOrderItem

	err := query.Where("sales_order_id = ?", PackingOrderId).Find(&res).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get PackingOrderItem list")
		return res, err
	}

	return res, nil
}

// GetPackingOrderItemById ...
func (d *postgreDatabase) GetPackingOrderItemById(ctx context.Context, id int64) (data models.PackingOrderItem, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("id = ?", id).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetPackingOrderByProductId ...
func (d *postgreDatabase) GetPackingOrderItemByProductId(ctx context.Context, productId int64) (data []models.PackingOrderItem, err error) {
	query := d.Db.WithContext(ctx)

    err = query.Where("product_id = ? ", productId).Preload("Customer").Preload("Product").First(&data).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


