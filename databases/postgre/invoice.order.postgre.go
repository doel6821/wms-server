package postgre

import (
	"context"
	"wms-server/constants"
	"wms-server/databases/postgre/models"

	"gorm.io/gorm"
)

// SaveInvoiceOrder ...
func (d *postgreDatabase) SaveInvoiceOrder(ctx context.Context, data models.InvoiceOrder) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveInvoiceOrder")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save InvoiceOrder")
		return err
	}
	return nil
}


// GetList ...
func (d *postgreDatabase) GetInvoiceOrderList(ctx context.Context, tenant, invoiceId int64, page , limit int) ( []models.InvoiceOrder, int64, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.InvoiceOrder
	var total int64

	query = query.Where("tenant = ?", tenant)
	if invoiceId != 0 {
		query = query.Where("customer_id ?", invoiceId)
	}

	err := query.Order("id desc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(0).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get InvoiceOrder list")
		return res, 0, err
	}

	return res, total, nil
}

// GetInvoiceOrderById ...
func (d *postgreDatabase) GetInvoiceOrderById(ctx context.Context, id int64) (data models.InvoiceOrderResponse, err error) {
	query := d.Db.WithContext(ctx)

	var order models.InvoiceOrderResponse
    err = query.Preload("InvoiceItems").First(&order, id).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


func (d *postgreDatabase) TxInvoiceOrder(ctx context.Context, reqInvoiceOrder models.InvoiceOrder, reqInvoiceOrderItem []models.InvoiceOrderItem ) error {
	query := d.Db.WithContext(ctx).Begin()

	query.SavePoint(constants.START)

	if err := d.TxSaveInvoiceOrder(ctx ,query, reqInvoiceOrder); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxSaveInvoiceOrderItems(ctx ,query, reqInvoiceOrderItem); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	return query.Commit().Error
}


func (d *postgreDatabase) TxSaveInvoiceOrder(ctx context.Context, query *gorm.DB, reqInvoiceOrder models.InvoiceOrder) error {
	if reqInvoiceOrder.ID == 0 {
		if err := query.Create(&reqInvoiceOrder).Error; err != nil {
			return err
		}
	} else {
		if err := query.Model(&reqInvoiceOrder).Updates(reqInvoiceOrder).Error; err != nil {
			return err
		}
	}

	return nil
}

func (d *postgreDatabase) TxSaveInvoiceOrderItems(ctx context.Context, query *gorm.DB, reqInvoiceOrderItem []models.InvoiceOrderItem) error {
	for _, v := range reqInvoiceOrderItem {
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



