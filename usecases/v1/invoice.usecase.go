package usecases

import (
	"context"
	"strconv"
	"strings"
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
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
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

	for _, v := range packDetail.Items {
		// create Invoice order items
		invoiceItem := pModels.InvoiceOrderItem{
			ProductCode: v.ProductCode,
			ProductName: v.ProductName,
			Quantity:    v.PackingOrderQty,
			Price:       v.Product.HETPrice,
			Total:       float64(v.PackingOrderQty) * (v.Product.HETPrice * (1 - float64(customer.DiscountPercent/100))),
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
		if len(product.ProductLocations) > 0 {
			productLocation, err := u.DB.GetPostgre().GetProductLocationByLocationCode(ctx, v.ProductId, product.ProductLocations[0].LocationCode)
			if err != nil {
				return res
			}
			productLocation.Qtty -= v.PackingOrderQty
			productLocations = append(productLocations, productLocation)
		}

		// update sales order item packing to invoice
		salesOrderItem, err := u.DB.GetPostgre().GetSalesOrderItemBySalesOrderId(ctx, v.SalesOrderId, v.ProductId)
		if err != nil {
			return res
		}

		salesOrderItem.PackingOrderQty -= v.PackingOrderQty
		salesOrderItem.InvoiceOrderQty += v.PackingOrderQty
		salesOrderItems = append(salesOrderItems, salesOrderItem)
	}

	paymentStatus := "Lunas"
	dueDate := time.Now()
	invoiceNumber := "INV/SO" + dueDate.Format("2006/01/02/150405")
	if customer.TermOfPayment != "cash" {
		lt, _ := strconv.Atoi(strings.Split(customer.TermOfPayment, " ")[0])
		dueDate = dueDate.Add(time.Duration(lt*24) * time.Hour)
		paymentStatus = "Belum Lunas"
	}

	inv := pModels.InvoiceOrder{
		PackingOrderId: req.PackingOrderId,
		InvoiceNumber: invoiceNumber,
		CustomerId:     customer.ID,
		CustomerName:   customer.Name,
		Discount:       int(customer.DiscountPercent),
		InvoiceDate:    time.Now(),
		Status:         constants.COMPLETE,
		PaymentStatus:  paymentStatus,
		DueDate:        dueDate,
		Tenant:         tenant,
	}

	packOrder := pModels.PackingOrder{
		ID:          packDetail.ID,
		CustomerId:  packDetail.CustomerId,
		PackingDate: packDetail.PackingDate,
		Status:      constants.COMPLETE,
		Tenant:      packDetail.Tenant,
	}
	// transactional
	// create invoice + invoice items
	// update packing order
	// update sales order item
	// update product and qtty on location
	inv.InvoiceItems = invoiceItems
	err = u.DB.GetPostgre().TxInvoiceOrder(ctx, inv, invoiceItems, packOrder, salesOrderItems, products, productLocations)
	if err != nil {
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetInvoiceList(ctx context.Context, tenant string, customerId int64, page, limit int, dueDate string) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listInvoiceOrder, total, err := u.DB.GetPostgre().GetInvoiceList(ctx, tenant, customerId, page, limit, dueDate)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get InvoiceOrder list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listInvoiceOrder
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetInvoiceDetailById(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	InvoiceOrder, err := u.DB.GetPostgre().GetInvoiceOrderById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Invoice")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = InvoiceOrder
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) UpdateInvoice(ctx context.Context, tenant string, req pModels.InvoiceOrder) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	err := u.DB.GetPostgre().SaveInvoiceOrder(ctx, req)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed update ReceiveOrder")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

