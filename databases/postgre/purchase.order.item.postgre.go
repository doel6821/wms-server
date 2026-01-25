package postgre

import (
	"context"
	"wms-server/databases/postgre/models"
	"wms-server/helpers"
)

// SavePurchaseOrderItem ...
func (d *postgreDatabase) SavePurchaseOrderItem(ctx context.Context, data []*models.PurchaseOrderItem) error {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SavePurchaseOrderItem")
	query := d.Db.WithContext(ctx)

	if err := query.Save(data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save PurchaseOrderItem")
		return err
	}
	return nil
}

// GetList ...
func (d *postgreDatabase) GetPurchaseOrderItemList(ctx context.Context, PurchaseOrderId int64) ([]models.PurchaseOrderItem, error) {
	query := d.Db.WithContext(ctx)

	var res []models.PurchaseOrderItem

	err := query.Where("sales_order_id = ?", PurchaseOrderId).Find(&res).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get PurchaseOrderItem list")
		return res, err
	}

	return res, nil
}

// GetPurchaseOrderItemById ...
func (d *postgreDatabase) GetPurchaseOrderItemById(ctx context.Context, id int64) (data models.PurchaseOrderItem, err error) {
	query := d.Db.WithContext(ctx)

	if err = query.Where("id = ?", id).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetPurchaseOrderItemByAllocation ...
func (d *postgreDatabase) GetPurchaseOrderItemByAllocation(ctx context.Context, ids []int) (data []models.PurchaseOrderItem, err error) {
	query := d.Db.WithContext(ctx)

	if err = query.Where("id in (" + helpers.JoinInts(ids, ",") + ") and allocation_order_quantity > 0").Find(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}

// GetPurchaseOrderItemByPurchaseOrderId ...
func (d *postgreDatabase) GetPurchaseOrderItemByPurchaseOrderId(ctx context.Context, PurchaseOrderId, productId int64) (data models.PurchaseOrderItem, err error) {
	query := d.Db.WithContext(ctx)

	if err = query.Where("purchase_order_id = ? and product_id = ?", PurchaseOrderId, productId).First(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	return data, nil
}


func (d *postgreDatabase) GetAvailableReceiveOrder(ctx context.Context, tenant string, supplierId, productId int64) (data []models.PurchaseOrderItem, err error) {
	query := d.Db.WithContext(ctx)

	if productId != 0 {
		err = query.Where("product_id = ? and ((order_qty - receive_order_qty) > 0 OR (order_qty - stocked_order_qty) > 0)", productId).Preload("Supplier").Preload("Product").Find(&data).Error; 
		if err != nil {
			d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		}
		
	} else if supplierId !=0 {
		err = query.Where("supplier_id = ? and ((order_qty - receive_order_qty) > 0 OR (order_qty - stocked_order_qty) > 0)", supplierId).Preload("Supplier").Preload("Product").Find(&data).Error; 
		if err != nil {
			d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		}
		
	}

	
	return data, nil
}


func (d *postgreDatabase) GetPurchaseOrderTotal(ctx context.Context, tenant string, startDate, endDate string) (data models.PurchaseOrderItemTotal, err error) {
	query := d.Db.WithContext(ctx)
	
	if err = query.Table("purchase_order_items poi").Select(`count(poi.product_id) as total_item, sum(poi.order_qty) as total_order_qty , 
		sum(poi.receive_order_qty) as total_receive_qty, sum(poi.stocked_order_qty) as total_stocked_qty, sum(poi.total) as total_amount`).
		Joins("left join purchase_orders po on po.id = poi.purchase_order_id").
		Where("(po.order_date between ? and ? ) and po.tenant = ?", startDate, endDate,tenant).Find(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return data, err
	}
	
	return data, nil
}

