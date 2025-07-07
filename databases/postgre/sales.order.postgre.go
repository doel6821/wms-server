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
func (d *postgreDatabase) GetSalesOrderList(ctx context.Context, tenant, customerId int64, page, limit int) ([]models.SalesOrder, int64, error) {
	query := d.Db.WithContext(ctx)

	var res []models.SalesOrder
	var total int64

	query = query.Where("tenant = ?", tenant)
	if customerId != 0 {
		query = query.Where("customer_id ?", customerId)
	}

	err := query.Order("id desc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(0).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get SalesOrder list")
		return res, 0, err
	}

	return res, total, nil
}

// GetSalesOrderById ...
func (d *postgreDatabase) GetSalesOrderById(ctx context.Context, id int64) (data models.SalesOrderResponse, err error) {
	query := d.Db.WithContext(ctx)

	var order models.SalesOrderResponse
	err = query.Preload("OrderItems.Product").First(&order, id).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

func (d *postgreDatabase) TxSalesOrder(ctx context.Context, reqSalesOrder models.SalesOrder, reqSalesOrderItem []models.SalesOrderItem) error {
	query := d.Db.WithContext(ctx).Begin()

	query.SavePoint(constants.START)

	if err := d.TxSaveSalesOrder(ctx, query, reqSalesOrder); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxSaveSalesOrderItems(ctx, query, reqSalesOrderItem); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	return query.Commit().Error
}

func (d *postgreDatabase) TxSaveSalesOrder(ctx context.Context, query *gorm.DB, reqSalesOrder models.SalesOrder) error {
	if reqSalesOrder.ID == 0 {
		if err := query.Create(&reqSalesOrder).Error; err != nil {
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
			if err := query.Model(&v).Updates(v).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
