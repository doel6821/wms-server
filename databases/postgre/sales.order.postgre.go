package postgre

import (
	"context"
	"wms-server/constants"
	"wms-server/databases/postgre/models"

	"gorm.io/gorm"
)

// SaveSalesOrder ...
func (d *postgreDatabase) SaveSalesOrder(ctx context.Context, data models.SalesOrder) error {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveSalesOrder")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save SalesOrder")
		return err
	}
	return nil
}

// GetList ...
func (d *postgreDatabase) GetSalesOrderList(ctx context.Context, tenant, allocation string, customerId int64, page, limit int) ([]models.SalesOrder, int64, error) {
	query := d.Db.WithContext(ctx)

	var res []models.SalesOrder
	var total int64
	var err error
	query = query.Where("tenant = ?", tenant)
	if customerId != 0 {
		query = query.Where("customer_id = ?", customerId)
	}

	if allocation == "true" {
		err = query.Preload("Customer").Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Where("sales_order_items.allocation_order_qty > 0")
			}).Order("id desc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error
	} else {
		err = query.Preload("Customer").Preload("Items").Order("id desc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error
	}

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get SalesOrder list :", err)
		return res, 0, err
	}

	var filtered []models.SalesOrder
	for _, so := range res {
		if len(so.Items) > 0 {
			filtered = append(filtered, so)
		}
	}

	return filtered, total, nil
}

// GetSalesOrderById ...
func (d *postgreDatabase) GetSalesOrderById(ctx context.Context, id int64) (data models.SalesOrder, err error) {
	query := d.Db.WithContext(ctx)
	
	err = query.Where("id = ? ", id).Preload("Customer").Preload("Items").Preload("Items.Product").First(&data).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

func (d *postgreDatabase) TxSalesOrder(ctx context.Context, reqSalesOrder models.SalesOrder, reqSalesOrderItem []models.SalesOrderItem, reqProductItems []models.Product, reqDemand []models.Demand) error {
	query := d.Db.WithContext(ctx).Begin()

	query.SavePoint(constants.START)
	if err := d.TxSaveProductItems(ctx, query, reqProductItems) ; err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxSaveSalesOrder(ctx, query, reqSalesOrder); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	// if err := d.TxSaveSalesOrderItems(ctx, query, reqSalesOrderItem); err != nil {
	// 	query.RollbackTo(constants.START)
	// 	return err
	// }

	if err := d.TxSaveDemands(ctx, query, reqDemand); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	return query.Commit().Error
}

func (d *postgreDatabase) TxSaveSalesOrder(ctx context.Context, query *gorm.DB, reqSalesOrder models.SalesOrder) error {
	if reqSalesOrder.ID == 0 {
		if err := query.Create(&reqSalesOrder).Error; err != nil {
		// if err := query.Preload("Items").First(&reqSalesOrder, reqSalesOrder.ID).Error; err != nil {
			return err
		}
	} else {
		if err := query.Model(&reqSalesOrder).Updates(reqSalesOrder).Error; err != nil {
			return err
		}
	}

	return nil
}

func (d *postgreDatabase) TxSaveSalesOrderItems(ctx context.Context, query *gorm.DB, reqSalesOrderItem []models.SalesOrderItem) error {
	for _, v := range reqSalesOrderItem {
		if v.ID == 0 {
			if err := query.Create(&v).Error; err != nil {
				return err
			}
		} else {
			if err := query.Model(&v).Updates(map[string]interface{}{
					"back_order_qty":    v.BackOrderQty,
					"allocation_order_qty": v.AllocationOrderQty,
					"packing_order_qty": v.PackingOrderQty,
					"invoice_order_qty":     v.InvoiceOrderQty,
					}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *postgreDatabase) TxSaveProductItems(ctx context.Context, query *gorm.DB, reqProductItems []models.Product) error {
	for _, v := range reqProductItems {	
		if err := query.Model(&v).Updates(v).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *postgreDatabase) TxSaveDemands(ctx context.Context, query *gorm.DB, reqDemand []models.Demand) error {
	for _, v := range reqDemand {
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
