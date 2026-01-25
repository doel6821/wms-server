package usecases

import (
	"context"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
)

func (u *usecase) CreateStocked(ctx context.Context, tenant string, req cModels.StockedRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	// get receiveOrderitems
	receiveDetail, err := u.DB.GetPostgre().GetReceiveOrderById(ctx, req.ReceiveOrderId)
	if err != nil {
		return res
	}

	products := []pModels.Product{}
	productLocations := []pModels.ProductLocation{}
	purchaseOrderItems := []pModels.PurchaseOrderItem{}
	salesOrderItems := []pModels.SalesOrderItem{}
	for _, v := range req.Items {
		// update product add qtty on hand
		product, err := u.DB.GetPostgre().GetProductById(ctx, v.ProductId)
		if err != nil {
			return res
		}

		boSales, err := u.DB.GetPostgre().GetSalesOrderItemBackOrderByProductId(ctx, v.ProductId)
		if err != nil {
			return res
		}

		allocation := 0
		receive := v.ReceiveOrderQty
		if len(boSales) > 0 {
			for i := 0 ; i < len(boSales) ; i++ {
				if boSales[i].BackOrderQty < receive {
					boSales[i].AllocationOrderQty = boSales[i].BackOrderQty
					boSales[i].BackOrderQty = 0
					allocation += boSales[i].AllocationOrderQty
					receive -= boSales[i].AllocationOrderQty
				} else if boSales[i].BackOrderQty > receive {
					boSales[i].AllocationOrderQty = receive
					boSales[i].BackOrderQty -= receive
					allocation += boSales[i].AllocationOrderQty
					receive -= boSales[i].AllocationOrderQty
					break
				}
			} 
		}

		salesOrderItems = append(salesOrderItems, boSales...)

		avgCost := ( (float64(product.StockOnHand) * product.AvgPrice) + (float64(v.ReceiveOrderQty) * v.PurchasePrice )) / ( float64( product.StockOnHand + v.ReceiveOrderQty ))
		product.StockAllocation += allocation
		product.StockBackOrder -= allocation
		product.StockOnReceive -= v.ReceiveOrderQty
		product.StockOnHand += receive
		product.CostPrice = v.PurchasePrice
		product.AvgPrice = avgCost
		products = append(products, product)
		
		// update product location deduct qtty
		productLocation, err := u.DB.GetPostgre().GetProductLocationByLocationCode(ctx, v.ProductId, v.ProductLocation)
		if err != nil {
			u.Logs.WithContext(ctx).Infoln("Location not set will set to new location")
			loc , err := u.DB.GetPostgre().GetLocationByLocationCode(ctx, tenant, v.ProductLocation)
			if err != nil {
				u.Logs.WithContext(ctx).Errorln("Location not found : ", err)
			}
			productLocation.LocationId = loc.ID
			productLocation.LocationCode = v.ProductLocation
			productLocation.ProductId = v.ProductId
		}

		productLocation.Qtty += v.ReceiveOrderQty
		productLocations = append(productLocations, productLocation)

		// update purchase order item
		purchaseOrderItem, err := u.DB.GetPostgre().GetPurchaseOrderItemByPurchaseOrderId(ctx, v.PurchaseOrderId, v.ProductId)
		if err != nil {
			return res
		}

		//purchaseOrderItem.ReceiveOrderQty -= v.ReceiveOrderQty
		purchaseOrderItem.StockedOrderQty += v.ReceiveOrderQty
		purchaseOrderItems = append(purchaseOrderItems, purchaseOrderItem)
	}

	receiveOrder := pModels.ReceiveOrder{
		ID:          receiveDetail.ID,
		SupplierId:  receiveDetail.SupplierId,
		ReceiveDate: receiveDetail.ReceiveDate,
		Status:      constants.COMPLETE,
		Tenant:      receiveDetail.Tenant,
	}
	// transactional
	// create Stocked + Stocked items
	// update packing order
	// update sales order item
	// update product and qtty on location
	err = u.DB.GetPostgre().TxStockedOrder(ctx, receiveOrder, purchaseOrderItems, products, productLocations, req.Items, salesOrderItems)
	if err != nil {
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}
