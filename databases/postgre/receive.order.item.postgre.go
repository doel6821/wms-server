package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
)

// SaveReceiveOrderItem ...
func (d *postgreDatabase) SaveReceiveOrderItem(ctx context.Context, data []*models.ReceiveOrderItem) error {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveReceiveOrderItem")
	query := d.Db.WithContext(ctx)

	if err := query.Save(data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save ReceiveOrderItem")
		return err
	}
	return nil
}

// GetList ...
func (d *postgreDatabase) GetReceiveOrderItemList(ctx context.Context, receiveOrderId int64) ([]models.ReceiveOrderItem, error) {
	query := d.Db.WithContext(ctx)

	var res []models.ReceiveOrderItem

	err := query.Where("receive_order_id = ?", receiveOrderId).Find(&res).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get ReceiveOrderItem list")
		return res, err
	}

	return res, nil
}

// GetReceiveOrderItemById ...
func (d *postgreDatabase) GetReceiveOrderItemById(ctx context.Context, id int64) (data models.ReceiveOrderItem, err error) {
	query := d.Db.WithContext(ctx)

	if err = query.Where("id = ?", id).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


// GetReceiveOrderItemByProductId ...
func (d *postgreDatabase) GetReceiveOrderItemByProductId(ctx context.Context, id int64) (data []models.ReceiveOrderItem, err error) {
	query := d.Db.WithContext(ctx)

	if err = query.Table("receive_order_items roi").Joins("left join receive_orders ro on ro.id = roi.receive_order_id").Where("roi.product_id = ? and roi.receive_order_qty > 0 and ro.status = ?", id, "on process").Find(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}
