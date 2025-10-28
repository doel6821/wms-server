package postgre

import (
	"context"
	"wms-server/constants"
	"wms-server/databases/postgre/models"

	"gorm.io/gorm"
)

// SaveInvoiceOrder ...
func (d *postgreDatabase) SaveInvoiceOrder(ctx context.Context, data models.InvoiceOrder) error {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveInvoiceOrder")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save InvoiceOrder")
		return err
	}
	return nil
}

// GetList ...
func (d *postgreDatabase) GetInvoiceList(ctx context.Context, tenant string, customerId int64, page, limit int, dueDate string) ([]models.InvoiceOrder, int64, error) {
	query := d.Db.WithContext(ctx)

	var res []models.InvoiceOrder
	var total int64

	query = query.Where("tenant = ?", tenant)
	if customerId != 0 {
		query = query.Where("customer_id = ?", customerId)
	}

	if dueDate != "" {
		query = query.Where("due_date between ? and ? ", dueDate+" 00:00:00", dueDate+" 23:59:59")
	}

	err := query.Order("id desc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error

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

// transactional
// create invoice + invoice items
// update packing order item
// update product and qtty on location

func (d *postgreDatabase) TxInvoiceOrder(ctx context.Context, reqInvoiceOrder models.InvoiceOrder, reqInvoiceOrderItem []models.InvoiceOrderItem, reqPackingOrder models.PackingOrder, reqSalesOrderItem []models.SalesOrderItem, products []models.Product, productLocations []models.ProductLocation) error {
	query := d.Db.WithContext(ctx).Begin()

	query.SavePoint(constants.START)

	if err := d.TxSaveInvoiceOrder(ctx, query, reqInvoiceOrder); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	// if err := d.TxSaveInvoiceOrderItems(ctx ,query, reqInvoiceOrderItem); err != nil {
	// 	query.RollbackTo(constants.START)
	// 	return err
	// }

	if err := d.TxUpdatePackingOrder(ctx, query, reqPackingOrder); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxUpdateSalesOrderItems(ctx, query, reqSalesOrderItem); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxUpdateProducts(ctx, query, products); err != nil {
		query.RollbackTo(constants.START)
		return err
	}

	if err := d.TxUpdateProductLocations(ctx, query, productLocations); err != nil {
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

func (d *postgreDatabase) TxUpdatePackingOrder(ctx context.Context, query *gorm.DB, reqPackingOrder models.PackingOrder) error {
	if err := query.Model(&models.PackingOrder{}).Where("id = ?", reqPackingOrder.ID).Updates(reqPackingOrder).Error; err != nil {
		return err
	}

	return nil
}

func (d *postgreDatabase) TxUpdateStatusPackingOrder(ctx context.Context, query *gorm.DB, id int64, status string) error {
	err := query.Model(&models.PackingOrder{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error getting existing data")
		return err
	}
	return nil
}

func (d *postgreDatabase) TxUpdateSalesOrderItems(ctx context.Context, query *gorm.DB, reqSalesOrderItem []models.SalesOrderItem) error {
	for i := range reqSalesOrderItem {
		v := reqSalesOrderItem[i]
		if v.ID == 0 {
			if err := query.Create(&v).Error; err != nil {
				return err
			}
		} else {
			if err := query.Model(&v).Updates(map[string]interface{}{
				"order_qty":            v.OrderQty,
				"allocation_order_qty": v.AllocationOrderQty,
				"back_order_qty":       v.BackOrderQty,
				"packing_order_qty":    v.PackingOrderQty,
				"invoice_order_qty":    v.InvoiceOrderQty,
			}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *postgreDatabase) TxUpdateProducts(ctx context.Context, query *gorm.DB, products []models.Product) error {
	d.Logs.WithContext(ctx).WithField("data", products).Info("UpdateProduct")

	for i := range products {
		v := products[i] // ambil value langsung dari slice
		if err := query.Model(&v).Updates(map[string]interface{}{
			"stock_on_hand":     v.StockOnHand,
			"stock_allocation":  v.StockAllocation,
			"stock_back_order":  v.StockBackOrder,
			"stock_on_purchase": v.StockOnPurchase,
			"stock_on_receive":  v.StockOnReceive,
			"stock_packing":     v.StockPacking,
			"cost_price":        v.CostPrice,
			"avg_price":         v.AvgPrice,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func (d *postgreDatabase) TxUpdateProductLocations(ctx context.Context, query *gorm.DB, productLocations []models.ProductLocation) error {
	d.Logs.WithContext(ctx).WithField("data", productLocations).Info("UpdateProduct")
	var err error
	for i := range productLocations {
		v := productLocations[i]
		if v.ID == 0 {
			err = query.Save(&v).Error
		} else if err = query.Model(&v).Updates(map[string]interface{}{
			"qtty": v.Qtty,
		}).Error; err != nil {
			return err
		}
	}
	return err
}
