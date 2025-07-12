package usecases

import (
	"context"
	"time"
	// "time"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
)

func (u *usecase) CreateInvoice(ctx context.Context, tenant string, req cModels.InvoiceRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR),
	}

	packDetail, err := u.DB.GetPostgre().GetPackingOrderById(ctx, req.PackingOrderId)
	if err != nil {
		return res
	}
	customer, err := u.DB.GetPostgre().GetCustomerById(ctx, packDetail.CustomerId)
	if err != nil {
		return res
	}
	// get packingorderitems where packingOrderId 
	invoiceItems := []pModels.InvoiceOrderItem{}
	products := []pModels.Product{}
	productLocations := []pModels.ProductLocation{}
	salesOrderItems := []pModels.SalesOrderItem{}
	
	for _, v := range packDetail.PackingItems {
		// create Invoice order items
		invoiceItem := pModels.InvoiceOrderItem{
			ProductCode    : v.ProductCode ,
			ProductName    : v.ProductName,
			Quantity       : v.PackingOrderQty,
			Price          : v.ProductPrice,
			Total          : float64(v.PackingOrderQty) * v.ProductPrice,
		}
		
		invoiceItems = append(invoiceItems, invoiceItem)
		
		// update product deduct qtty on packing
		product, err := u.DB.GetPostgre().GetProductById(ctx, v.ProductId)
		if err != nil {
			return res
		}

		product.StockPacking -= v.PackingOrderQty
		products = append(products, product)

		// update product location deduct qtty 
		productLocation , err := u.DB.GetPostgre().GetProductLocationByLocationCode(ctx, v.ProductId, v.Location)
		if err != nil {
			return res
		}

		productLocation.Qtty -= v.PackingOrderQty
		productLocations = append(productLocations, productLocation)

		// update sales order item packing to invoice 
		salesOrderItem, err := u.DB.GetPostgre().GetSalesOrderItemBySalesOrderId(ctx, v.SalesOrderId, v.ProductId)
		if err != nil {
			return res
		}
		
		salesOrderItem.PackingOrderQty -= v.PackingOrderQty
		salesOrderItem.InvoiceOrderQty += v.PackingOrderQty
		salesOrderItems = append(salesOrderItems, salesOrderItem)
	}

	inv := pModels.InvoiceOrder{
		PackingOrderId : req.PackingOrderId,
		CustomerId     : customer.ID,
		CustomerName   : customer.Name,
		InvoiceDate    : time.Now(),
		Status         : constants.COMPLETE,
		Tenant         : tenant, 
	}

	packOrder := pModels.PackingOrder{
		ID          : packDetail.ID,
		CustomerId  : packDetail.CustomerId,
		PackingDate : packDetail.PackingDate,
		Status      : constants.COMPLETE,
		Tenant      : packDetail.Tenant,
	}
	// transactional
	// create invoice + invoice items
	// update packing order
	// update sales order item
	// update product and qtty on location
	
	err = u.DB.GetPostgre().TxInvoiceOrder(ctx, inv, invoiceItems, packOrder, salesOrderItems, products, productLocations)
	if err != nil {
		return res
	}

	res.Meta = helpers.GetNewMetaResponse("en", constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetInvoiceList(ctx context.Context,tenant string, customerId int64, page, limit int) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR),
	}

	listInvoiceOrder, total, err := u.DB.GetPostgre().GetInvoiceList(ctx, tenant, customerId, page, limit)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get InvoiceOrder list")
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listInvoiceOrder
	res.Count = total
	res.Meta = helpers.GetNewMetaResponse("en", constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetInvoiceDetailById(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR),
	}

	InvoiceOrder, err := u.DB.GetPostgre().GetInvoiceOrderById(ctx, id)
	if err != nil && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Product")
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = InvoiceOrder
	res.Meta = helpers.GetNewMetaResponse("en", constants.RC_SUCCESS)
	return res
}


