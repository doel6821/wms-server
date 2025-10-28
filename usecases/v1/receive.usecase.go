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

func (u *usecase) CreateReceive(ctx context.Context, tenant string, req cModels.ReceiveOrderRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	supplier, err := u.DB.GetPostgre().GetSupplierById(ctx, req.SupplierID)
	if err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	paymentStatus := "Lunas"
	dueDate := time.Now()
	if supplier.TermOfPayment != "cash" {
		lt, _ := strconv.Atoi(strings.Split(supplier.TermOfPayment, " ")[0])
		dueDate = dueDate.Add(time.Duration(lt*24) * time.Hour)
		paymentStatus = "Belum Lunas"
	}

	receive := pModels.ReceiveOrder{
		SupplierId:    req.SupplierID,
		ReceiveDate:   time.Now(),
		Status:        constants.ON_PROCESS,
		InvoiceNumber: req.InvoiceNumber,
		PaymentStatus: paymentStatus,
		DueDate:       dueDate,
		Tenant:        tenant,
	}

	receiveItems := []pModels.ReceiveOrderItem{}
	products := []pModels.Product{}
	purchaseOrderItems := []pModels.PurchaseOrderItem{}

	for _, v := range req.ReceiveOrders {
		product, err := u.DB.GetPostgre().GetProductById(ctx, v.ProductId)
		if err != nil {
			res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
			return res
		}
		location := ""
		if len(product.ProductLocations) > 0 {
			location = product.ProductLocations[0].LocationCode
		}
		// create Receive order items
		receiveItem := pModels.ReceiveOrderItem{
			PurchaseOrderId: v.PurchaseOrderID,
			ProductId:       v.ProductId,
			ProductCode:     product.Code,
			ProductName:     product.Name,
			ProductLocation: location,
			ReceiveOrderQty: v.ReceiveOrderQtty,
			PurchasePrice:   v.PurchasePrice,
		}

		receiveItems = append(receiveItems, receiveItem)

		// update product deduct qtty on packing

		product.StockOnPurchase -= v.ReceiveOrderQtty
		product.StockOnReceive += v.ReceiveOrderQtty
		products = append(products, product)

		// update purchase order item order to Receive
		purchaseOrderItem, err := u.DB.GetPostgre().GetPurchaseOrderItemByPurchaseOrderId(ctx, v.PurchaseOrderID, v.ProductId)
		if err != nil {
			return res
		}

		purchaseOrderItem.ReceiveOrderQty += v.ReceiveOrderQtty
		purchaseOrderItems = append(purchaseOrderItems, purchaseOrderItem)

	}
	// transactional
	// create Receive + Receive items
	// update purchase order item
	// update product and qtty on receive
	receive.Items = receiveItems

	err = u.DB.GetPostgre().TxReceiveOrder(ctx, receive, receiveItems, purchaseOrderItems, products)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed create ReceiveOrder")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetReceiveList(ctx context.Context, tenant string, customerId int64, page, limit int, dueDate string) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listReceiveOrder, total, err := u.DB.GetPostgre().GetReceiveOrderList(ctx, tenant, customerId, page, limit, dueDate)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get ReceiveOrder list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listReceiveOrder
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetReceiveDetailById(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	ReceiveOrder, err := u.DB.GetPostgre().GetReceiveOrderById(ctx, id)
	if err != nil && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Product")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = ReceiveOrder
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetReceiveDetailByProductId(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	ReceiveOrder, err := u.DB.GetPostgre().GetReceiveOrderItemByProductId(ctx, id)
	if err != nil && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Product")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = ReceiveOrder
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetAvailableReceiveList(ctx context.Context, tenant string, supplierId, productId int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	availableReceiveOrder, err := u.DB.GetPostgre().GetAvailableReceiveOrder(ctx, tenant, supplierId, productId)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get AvailableReceiveOrder list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = availableReceiveOrder
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) UpdateReceive(ctx context.Context, tenant string, req pModels.ReceiveOrder) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	err := u.DB.GetPostgre().SaveReceiveOrder(ctx, req)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed update ReceiveOrder")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}
