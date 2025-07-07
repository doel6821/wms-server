package postgre

import (
	"context"
	"wms-server/constants"
	"wms-server/databases/postgre/models"

	"gorm.io/gorm"
)

// SavePackingOrder ...
func (d *postgreDatabase) SavePackingOrder(ctx context.Context, data models.PackingOrder) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SavePackingOrder")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save PackingOrder")
		return err
	}
	return nil
}


// GetList ...
func (d *postgreDatabase) GetPackingOrderList(ctx context.Context, tenant, customerId int64, page , limit int) ( []models.PackingOrder, int64, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.PackingOrder
	var total int64

	query = query.Where("tenant = ?", tenant)
	if customerId != 0 {
		query = query.Where("customer_id ?", customerId)
	}

	err := query.Order("id desc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(0).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get PackingOrder list")
		return res, 0, err
	}

	return res, total, nil
}

// GetPackingOrderById ...
func (d *postgreDatabase) GetPackingOrderById(ctx context.Context, id int64) (data models.PackingOrderResponse, err error) {
	query := d.Db.WithContext(ctx)

	var order models.PackingOrderResponse
    err = query.Preload("PackingItems").First(&order, id).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


func (d *postgreDatabase) TxPackingOrder(ctx context.Context, reqPackingOrder models.PackingOrder, reqPackingOrderItem []models.PackingOrderItem ) error {
	query := d.Db.WithContext(ctx).Begin()

	query.SavePoint(constants.START)

	if err := d.TxSavePackingOrder(ctx ,query, reqPackingOrder); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxSavePackingOrderItems(ctx ,query, reqPackingOrderItem); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	return query.Commit().Error
}


func (d *postgreDatabase) TxSavePackingOrder(ctx context.Context, query *gorm.DB, reqPackingOrder models.PackingOrder) error {
	if reqPackingOrder.ID == 0 {
		if err := query.Create(&reqPackingOrder).Error; err != nil {
			return err
		}
	} else {
		if err := query.Model(&reqPackingOrder).Updates(reqPackingOrder).Error; err != nil {
			return err
		}
	}

	return nil
}

func (d *postgreDatabase) TxSavePackingOrderItems(ctx context.Context, query *gorm.DB, reqPackingOrderItem []models.PackingOrderItem) error {
	for _, v := range reqPackingOrderItem {
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



