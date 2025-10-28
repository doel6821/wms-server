package usecases

import (
	"context"
	"time"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
)

func (u *usecase) CreatePurchaseOrder(ctx context.Context, tenant string, req cModels.PurchaseOrderRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	supplier, err := u.DB.GetPostgre().GetSupplierById(ctx, req.SupplierID)
	if err != nil && err.Error() != "record not found" {
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	} else if supplier.ID == 0 {
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR) // custNotFound
		return res
	}
	purchaseNumber := "PO/" + time.Now().Format("2006/01/02/150405")
	po := pModels.PurchaseOrder{
		PurchaseNumber: purchaseNumber,
		SupplierId:     req.SupplierID,
		OrderDate:      time.Now(),
		Amount:         req.Amount,
		Discount:       req.Discount,
		Total:          req.TotalAmount,
		Tenant:         tenant,
	}

	poItems := []pModels.PurchaseOrderItem{}
	productItems := []pModels.Product{}

	for _, v := range req.OrderItems {
		// checkStock
		product, err := u.DB.GetPostgre().GetProductById(ctx, int64(v.ProductId))
		if err != nil {
			res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
			return res
		}

		orderItem := pModels.PurchaseOrderItem{}
		orderItem.SupplierId = product.SupplierId
		orderItem.ProductId = product.ID
		orderItem.Price = product.HETPrice
		orderItem.OrderQty = v.OrderQty
		orderItem.SubTotal = v.Price * float64(orderItem.OrderQty)
		orderItem.Discount = v.Price * float64(req.Discount/100)
		orderItem.Total = orderItem.SubTotal - orderItem.Discount

		product.StockOnPurchase += v.OrderQty
		productItems = append(productItems, product)
		poItems = append(poItems, orderItem)

	}
	po.Items = poItems
	err = u.DB.GetPostgre().TxPurchaseOrder(ctx, po, poItems, productItems)
	if err != nil {
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetPurchaseOrderList(ctx context.Context, tenant string, customerId int64, page, limit int) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listPurchaseOrder, total, err := u.DB.GetPostgre().GetPurchaseOrderList(ctx, tenant, customerId, page, limit)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get PurchaseOrder list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listPurchaseOrder
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetPurchaseOrderDetailById(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	PurchaseOrder, err := u.DB.GetPostgre().GetPurchaseOrderById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Product")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = PurchaseOrder
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}
