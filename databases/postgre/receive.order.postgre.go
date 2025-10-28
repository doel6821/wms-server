package postgre

import (
	"context"
	"wms-server/constants"
	"wms-server/databases/postgre/models"

	"gorm.io/gorm"
)

// SaveReceiveOrder ...
func (d *postgreDatabase) SaveReceiveOrder(ctx context.Context, data models.ReceiveOrder) error {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveReceiveOrder")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save ReceiveOrder")
		return err
	}
	return nil
}

// GetList ...
func (d *postgreDatabase) GetReceiveOrderList(ctx context.Context, tenant string, supplierId int64, page, limit int, dueDate string) ([]models.ReceiveOrder, int64, error) {
	query := d.Db.WithContext(ctx)

	var res []models.ReceiveOrder
	var total int64

	query = query.Where("tenant = ?", tenant)
	if supplierId != 0 {
		query = query.Where("supplier_id ?", supplierId)
	}

	if dueDate != "" {
		query = query.Where("due_date between ? and ? ", dueDate+" 00:00:00", dueDate+" 23:59:59")
	}

	err := query.Preload("Supplier").Preload("Items").Order("id desc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get ReceiveOrder list")
		return res, 0, err
	}

	return res, total, nil
}

// GetReceiveOrderById ...
func (d *postgreDatabase) GetReceiveOrderById(ctx context.Context, id int64) (data models.ReceiveOrder, err error) {
	query := d.Db.WithContext(ctx)

	err = query.Preload("Supplier").Preload("Items").First(&data, id).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetReceiveOrderByProductId ...
func (d *postgreDatabase) GetReceiveOrderByProductId(ctx context.Context, id int64) (data []models.ReceiveOrder, err error) {
	query := d.Db.WithContext(ctx)

	err = query.Preload("Supplier").Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Where("receive_order_items.product_id = ? and receive_order_items.receive_order_qty > 0", id)
	}).Find(&data).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// UpdateReceiveOrderById ...
func (d *postgreDatabase) UpdateReceiveOrder(ctx context.Context, data models.ReceiveOrder) error {
	query := d.Db.WithContext(ctx)

	err := query.Updates(&data).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return err
	}
	return nil
}

// UpdateReceiveOrderById ...
func (d *postgreDatabase) UpdateStatusReceiveOrder(ctx context.Context, id int64, status string) error {
	query := d.Db.WithContext(ctx)

	err := query.Model(&models.ReceiveOrder{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return err
	}
	return nil
}

// transactional
// create Receive + Receive items
// update purchase order item
// update product and qtty on receive
func (d *postgreDatabase) TxReceiveOrder(ctx context.Context, reqReceiveOrder models.ReceiveOrder, reqReceiveOrderItem []models.ReceiveOrderItem, reqPurchaseOrderItems []models.PurchaseOrderItem, reqProducts []models.Product) error {
	query := d.Db.WithContext(ctx).Begin()

	query.SavePoint(constants.START)

	if err := d.TxSaveReceiveOrder(ctx, query, reqReceiveOrder, reqReceiveOrderItem); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	// if err := d.TxSaveReceiveOrderItems(ctx ,query, reqReceiveOrderItem); err != nil {
	// 	query.RollbackTo(constants.START)
	// 	return err
	// }

	if err := d.TxUpdatePurchaseOrderItems(ctx, query, reqPurchaseOrderItems); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxUpdateProducts(ctx, query, reqProducts); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	return query.Commit().Error
}

func (d *postgreDatabase) TxStockedOrder(ctx context.Context, reqReceiveOrder models.ReceiveOrder, reqPurchaseOrderItems []models.PurchaseOrderItem, reqProducts []models.Product, reqProductLocations []models.ProductLocation, reqReceiveOrderItems []models.ReceiveOrderItem, boSales []models.SalesOrderItem) error {
	query := d.Db.WithContext(ctx).Begin()

	query.SavePoint(constants.START)

	if err := d.TxSaveReceiveOrder(ctx, query, reqReceiveOrder, reqReceiveOrderItems); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxUpdatePurchaseOrderItems(ctx, query, reqPurchaseOrderItems); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxUpdateProducts(ctx, query, reqProducts); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxUpdateProductLocations(ctx, query, reqProductLocations); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxSaveSalesOrderItems(ctx, query, boSales); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	

	return query.Commit().Error
}

func (d *postgreDatabase) TxSaveReceiveOrder(ctx context.Context, query *gorm.DB, reqReceiveOrder models.ReceiveOrder, reqReceiveOrderItems []models.ReceiveOrderItem) error {
	if reqReceiveOrder.ID == 0 {
		if err := query.Create(&reqReceiveOrder).Error; err != nil {
			return err
		}
	} else {
		if err := query.Model(&reqReceiveOrder).Updates(reqReceiveOrder).Error; err != nil {
			return err
		}

		for _, item := range reqReceiveOrderItems {
			if err := query.Model(&models.ReceiveOrderItem{}).
				Where("id = ?", item.ID).
				Updates(map[string]interface{}{
					"product_location": item.ProductLocation,
				}).Error; err != nil {
				return err
			}
		}

	}

	return nil
}

func (d *postgreDatabase) TxUpdatePurchaseOrderItems(ctx context.Context, query *gorm.DB, reqPurchaseOrderItem []models.PurchaseOrderItem) error {
	for _, v := range reqPurchaseOrderItem {
		if err := query.Model(&v).Updates(v).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *postgreDatabase) TxSaveReceiveOrderItems(ctx context.Context, query *gorm.DB, reqReceiveOrderItem []models.ReceiveOrderItem) error {
	for _, v := range reqReceiveOrderItem {
		if v.ID == 0 {
			if err := query.Create(&v).Error; err != nil {
				return err
			}
		} else {
			if err := query.Model(&v).Updates(v).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
