package postgre

import (
	"context"
	"wms-server/constants"
	"wms-server/databases/postgre/models"

	"gorm.io/gorm"
)

// SavePurchaseOrder ...
func (d *postgreDatabase) SavePurchaseOrder(ctx context.Context, data models.PurchaseOrder) error {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SavePurchaseOrder")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save PurchaseOrder")
		return err
	}
	return nil
}

// GetList ...
func (d *postgreDatabase) GetPurchaseOrderList(ctx context.Context, tenant string, supplierId int64, page, limit int) ([]models.PurchaseOrder, int64, error) {
	query := d.Db.WithContext(ctx)

	var res []models.PurchaseOrder
	var total int64

	query = query.Where("tenant = ?", tenant)
	if supplierId != 0 {
		query = query.Where("supplier_id ?", supplierId)
	}

	err := query.Preload("Supplier").Preload("Items").Order("id desc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get PurchaseOrder list")
		return res, 0, err
	}

	return res, total, nil
}

// GetPurchaseOrderById ...
func (d *postgreDatabase) GetPurchaseOrderById(ctx context.Context, id int64) (data models.PurchaseOrder, err error) {
	query := d.Db.WithContext(ctx)

	err = query.Where("id = ? ", id).Preload("Supplier").Preload("Items").Preload("Items.Product").First(&data).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

func (d *postgreDatabase) TxPurchaseOrder(ctx context.Context, reqPurchaseOrder models.PurchaseOrder, reqPurchaseOrderItem []models.PurchaseOrderItem, reqProductItems []models.Product) error {
	query := d.Db.WithContext(ctx).Begin()

	query.SavePoint(constants.START)
	if err := d.TxSaveProductItems(ctx, query, reqProductItems) ; err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxSavePurchaseOrder(ctx, query, reqPurchaseOrder); err != nil {
		query.RollbackTo(constants.START)
		return err
	}


	return query.Commit().Error
}

func (d *postgreDatabase) TxSavePurchaseOrder(ctx context.Context, query *gorm.DB, reqPurchaseOrder models.PurchaseOrder) error {
	
	if err := query.Create(&reqPurchaseOrder).Error; err != nil {
		return err
	}


	return nil
}

func (d *postgreDatabase) TxSavePurchaseOrderItems(ctx context.Context, query *gorm.DB, reqPurchaseOrderItem []models.PurchaseOrderItem) error {
	for _, v := range reqPurchaseOrderItem {
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
