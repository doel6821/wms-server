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

func (u *usecase) CreateSalesOrder(ctx context.Context, tenant string, req cModels.SalesOrderRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR),
	}

	cust, err := u.DB.GetPostgre().GetCustomerById(ctx, req.CustomerID)
	if err != nil && err.Error() != "record not found" {
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR)
		return res
	} else if cust.ID == 0 {
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR) // custNotFound
		return res
	}

	so := pModels.SalesOrder{
		CustomerId : uint(req.CustomerID) ,
		OrderDate  : time.Now() ,
		Amount     : req.Amount ,
		Discount   : req.Discount ,
		Total      : req.TotalAmount ,
		Tenant     : tenant ,
	}

	soItems := []pModels.SalesOrderItem{}
	productItems := []pModels.Product{}
	demands := []pModels.Demand{}
	for _, v := range req.OrderItems {
		// checkStock
		product , err := u.DB.GetPostgre().GetProductById(ctx, int64(v.ProductId))
		if err != nil && err.Error() != "record not found" {
			res.Meta = helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR)
			return res
		}

		orderItem := pModels.SalesOrderItem{}
		allocation := 0
		backOrder := 0
		// update stock to allocation if ready stock to backOrder if no stock
		if product.StockOnHand != 0 && product.StockOnHand > v.OrderQty {
			allocation = v.OrderQty
			
		} else if product.StockOnHand != 0 && product.StockOnHand < v.OrderQty {
			allocation = product.StockOnHand
			backOrder = v.OrderQty - allocation

		} else if product.StockOnHand == 0 {
			backOrder = v.OrderQty
		}
		product.StockOnHand -= allocation
		product.StockAllocation += allocation
		product.StockBackOrder += backOrder
		
		orderItem.ProductId = product.ID
		orderItem.Price = product.HETPrice
		orderItem.AllocationOrderQty = allocation
		orderItem.BackOrderQty = backOrder
		orderItem.Total = orderItem.Price * float64(orderItem.AllocationOrderQty)
		
		productItems = append(productItems, product)
		soItems = append(soItems, orderItem)

		// update demand
		demand, err := u.DB.GetPostgre().GetDemandByProductId(ctx, int64(v.ProductId))
		if err != nil && err.Error() != "record not found" {
			res.Meta = helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR)
			return res
		}
		trxDate := time.Now().Format("2006-01")
		if demand.ID == 0 {
			demand.Month = trxDate
			demand.NQty = int64(v.OrderQty)
			demand.ProductId = int64(v.ProductId)
		} else {
			if trxDate != demand.Month {
				demand.Month = trxDate
				demand.N12Qty = demand.N11Qty
				demand.N11Qty = demand.N10Qty
				demand.N10Qty = demand.N9Qty
				demand.N9Qty = demand.N8Qty
				demand.N8Qty = demand.N7Qty
				demand.N7Qty = demand.N6Qty
				demand.N6Qty = demand.N5Qty
				demand.N5Qty = demand.N4Qty
				demand.N4Qty = demand.N3Qty
				demand.N3Qty = demand.N2Qty
				demand.N2Qty = demand.N1Qty
				demand.N1Qty = demand.NQty
				demand.NQty = int64(v.OrderQty)
			} else {
				demand.NQty += int64(v.OrderQty)
			}
		}
		demands = append(demands, demand)
	}

	
	err = u.DB.GetPostgre().TxSalesOrder(ctx, so, soItems, productItems, demands)
	if err != nil {
		return res
	}

	res.Meta = helpers.GetNewMetaResponse("en", constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetSalesOrderList(ctx context.Context,tenant string, customerId int64, page, limit int) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR),
	}

	listSalesOrder, total, err := u.DB.GetPostgre().GetSalesOrderList(ctx, tenant, customerId, page, limit)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get SalesOrder list")
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listSalesOrder
	res.Count = total
	res.Meta = helpers.GetNewMetaResponse("en", constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetSalesOrderDetailById(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR),
	}

	salesOrder, err := u.DB.GetPostgre().GetSalesOrderById(ctx, id)
	if err != nil && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Product")
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = salesOrder
	res.Meta = helpers.GetNewMetaResponse("en", constants.RC_SUCCESS)
	return res
}


