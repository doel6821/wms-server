package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
)

// SaveInvoiceOrderItem ...
func (d *postgreDatabase) SaveInvoiceOrderItem(ctx context.Context, data []*models.InvoiceOrderItem) (error) {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveInvoiceOrderItem")
	query := d.Db.WithContext(ctx)

	if err := query.Save(data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save InvoiceOrderItem")
		return err
	}
	return nil
}


// GetList ...
func (d *postgreDatabase) GetInvoiceOrderItemList(ctx context.Context, InvoiceOrderId int64) ( []models.InvoiceOrderItem, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.InvoiceOrderItem

	err := query.Where("invoice_order_id = ?", InvoiceOrderId).Find(&res).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get InvoiceOrderItem list")
		return res, err
	}

	return res, nil
}

// GetInvoiceOrderItemById ...
func (d *postgreDatabase) GetInvoiceOrderItemById(ctx context.Context, id int64) (data models.InvoiceOrderItem, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Where("id = ?", id).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


