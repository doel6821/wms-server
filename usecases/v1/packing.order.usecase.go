package usecases

import (
	"context"
	"fmt"
	"time"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
)

func (u *usecase) CreatePackingOrder(ctx context.Context, tenant string, req cModels.PackingOrderRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	pack := pModels.PackingOrder{
		CustomerId:  req.CustomerID,
		PackingDate: time.Now(),
		Status:      constants.ON_PACKING,
		Tenant:      tenant,
	}

	// get salesorderitems where allocation > 0
	allocationItems, err := u.DB.GetPostgre().GetSalesOrderItemByAllocation(ctx, req.SalesOrderIDs)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Sales Order Item Allocation list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	salesOrderItems := []pModels.SalesOrderItem{}
	packingItems := []pModels.PackingOrderItem{}
	products := []pModels.Product{}
	for _, v := range allocationItems {
		product, err := u.DB.GetPostgre().GetProductById(ctx, v.ProductId)
		if err != nil {
			u.Logs.WithContext(ctx).WithError(err).Error("failed get Product Location list")
			res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
			return res
		}

		if len(product.ProductLocations) > 0 {
			if product.ProductLocations[0].Qtty >= v.AllocationOrderQty {
				packingItem := pModels.PackingOrderItem{
					SalesOrderId:    v.SalesOrderId,
					CustomerId:      req.CustomerID,
					ProductId:       v.ProductId,
					ProductCode:     product.Code,
					ProductName:     product.Name,
					ProductLocation: product.ProductLocations[0].LocationCode,
					PackingOrderQty: v.AllocationOrderQty,
				}
				packingItems = append(packingItems, packingItem)
			} else if product.ProductLocations[0].Qtty < v.AllocationOrderQty && product.ProductLocations[0].Qtty+product.ProductLocations[1].Qtty >= v.AllocationOrderQty {
				packingItem := pModels.PackingOrderItem{
					SalesOrderId:    v.SalesOrderId,
					CustomerId:      req.CustomerID,
					ProductId:       v.ProductId,
					ProductCode:     product.Code,
					ProductName:     product.Name,
					ProductLocation: product.ProductLocations[0].LocationCode,
					PackingOrderQty: product.ProductLocations[0].Qtty,
				}
				packingItems = append(packingItems, packingItem)
	
				packingItem = pModels.PackingOrderItem{
					SalesOrderId:    v.SalesOrderId,
					CustomerId:      req.CustomerID,
					ProductId:       v.ProductId,
					ProductCode:     product.Code,
					ProductName:     product.Name,
					ProductLocation: product.ProductLocations[1].LocationCode,
					PackingOrderQty: v.AllocationOrderQty - product.ProductLocations[0].Qtty,
				}
				packingItems = append(packingItems, packingItem)
			}
		} else {
			packingItem := pModels.PackingOrderItem{
					SalesOrderId:    v.SalesOrderId,
					CustomerId:      req.CustomerID,
					ProductId:       v.ProductId,
					ProductCode:     product.Code,
					ProductName:     product.Name,
					ProductLocation: "-",
					PackingOrderQty: v.AllocationOrderQty,
				}
				packingItems = append(packingItems, packingItem)
		}

		// update allocation to packing
		fmt.Println(v.AllocationOrderQty, "v.alokasi")
		fmt.Println(product.StockAllocation, "product.alokasi awal")
		product.StockAllocation -= v.AllocationOrderQty
		product.StockPacking += v.AllocationOrderQty
		fmt.Println(product.StockAllocation, "product.alokasi akhir")
		products = append(products, product)

		// update sales order item allocation to pick
		salesOrderItem := pModels.SalesOrderItem{
			ID:                 v.ID,
			SalesOrderId:       v.SalesOrderId,
			CustomerId:         req.CustomerID,
			ProductId:          v.ProductId,
			OrderQty:           v.OrderQty,
			BackOrderQty:       v.BackOrderQty,
			AllocationOrderQty: 0,
			PackingOrderQty:    v.AllocationOrderQty,
			InvoiceOrderQty:    v.InvoiceOrderQty,
			Price:              v.Price,
			Total:              v.Total,
		}
		salesOrderItems = append(salesOrderItems, salesOrderItem)
	}
	pack.Items = packingItems
	// save packing order & packing order items
	// update sales order item allocation to pick
	// update product allocation to packing
	err = u.DB.GetPostgre().TxPackingOrder(ctx, pack, packingItems, salesOrderItems, products)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed create Packing Order")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetPackingOrderList(ctx context.Context, tenant string, customerId int64, page, limit int) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listPackingOrder, total, err := u.DB.GetPostgre().GetPackingOrderList(ctx, tenant, customerId, page, limit)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get PackingOrder list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listPackingOrder
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetPackingOrderDetailById(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	PackingOrder, err := u.DB.GetPostgre().GetPackingOrderById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Packing Order Detail")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = PackingOrder
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}


func (u *usecase) GetPackingOrderListByProductId(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	PackingOrder, err := u.DB.GetPostgre().GetPackingOrderItemList(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Packing Order Detail")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = PackingOrder
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}


