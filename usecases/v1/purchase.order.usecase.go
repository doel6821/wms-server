package usecases

import (
	"context"
	"fmt"
	"math"
	"strconv"
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


func (u *usecase) GetPurchaseOrderRecomendation(ctx context.Context, tenant string,  supplierId int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	configMinStock, err := u.DB.GetPostgre().GetConfigurationByConfigurationName(ctx, tenant, "MINIMUM_STOCK")
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get config")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	configBufferStock, err := u.DB.GetPostgre().GetConfigurationByConfigurationName(ctx, tenant, "BUFFER_STOCK")
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get config")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	minStockParam, err := strconv.ParseFloat(configMinStock.Value, 64)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed parse float")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	fmt.Println("mintockParam : ", minStockParam)
	bufferStockParam, err := strconv.ParseFloat(configBufferStock.Value, 64)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed parse float")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	fmt.Println("bufferStockParam : ", bufferStockParam)
	products, _ , err := u.DB.GetPostgre().GetProductList(ctx, "", tenant, "", strconv.Itoa(int(supplierId)), 1, 99999)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Product")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	recomend := []pModels.PurchaseOrderRecomendation{}

	for _, product := range products {
		if product.Demands.ID != 0 {
			fmt.Println("========================")
			totalDemand12Months := product.Demands.NQty +
				product.Demands.N1Qty +
				product.Demands.N2Qty +
				product.Demands.N3Qty +
				product.Demands.N4Qty +
				product.Demands.N5Qty +
				product.Demands.N6Qty +
				product.Demands.N7Qty +
				product.Demands.N8Qty +
				product.Demands.N9Qty +
				product.Demands.N10Qty +
				product.Demands.N11Qty 
				// product.Demands.N12Qty
			
			if totalDemand12Months == 0 {
				// Jika tidak ada demand, lewati (tidak perlu dibeli)
				continue 
			}

			age := helpers.MonthsElapsed(product.CreatedAt, time.Now())
			if age >= 12 {
				age = 12
			} else if age == 0 {
				age = 1
			}

			fmt.Println("product createdAt : " ,product.CreatedAt)
			fmt.Println("product age : ", age)
			avgDemandMonthly := totalDemand12Months / int64(age)
			fmt.Println("average demand : ", avgDemandMonthly)
			minStock := float64(avgDemandMonthly) * minStockParam
			fmt.Println("minimum stock : ", minStock)
			bufferStock := float64(avgDemandMonthly) * (bufferStockParam/30)
			fmt.Println( "buffer stock : ", bufferStock )
			futureStock := product.StockOnHand + product.StockOnPurchase + product.StockOnReceive
			fmt.Println( "future stock : ", futureStock )

			if float64(futureStock) > (minStock + bufferStock) {
				u.Logs.WithContext(ctx).Printf("SKIP: Product %s (ID: %d). FutureStock (%d) > Filter Threshold (%.2f)", product.Name, product.ID, futureStock, minStock)
				continue // Lanjutkan ke produk berikutnya
			}
			
			projectedDemand := float64(avgDemandMonthly) * (float64(product.Supplier.LeadTimeDays)/30)
			fmt.Println("projected demand : ", projectedDemand)
			projectedEndStock := float64(futureStock) - projectedDemand
			fmt.Println("projected end stock : ", projectedEndStock)
			if projectedEndStock < (minStock + bufferStock) {
				targetStock := minStock + bufferStock
				fmt.Println("target stock : ", targetStock)
				purchaseQuantityFloat := targetStock - projectedEndStock

				purchaseQuantity := int(math.Ceil(math.Max(0, purchaseQuantityFloat)))
				fmt.Println("purchase Qty : ", purchaseQuantity)
				u.Logs.WithContext(ctx).Printf("RECOMMEND: Product %s. End Stock: %.2f vs Safety Stock: %.2f. PR Qty: %d", product.Name, projectedEndStock, minStock, purchaseQuantity)

				purchaseProduct := pModels.PurchaseOrderRecomendation{
					SupplierId: product.SupplierId,
					ProductId: product.ID,
					ProductName: product.Name,
					OrderQty: purchaseQuantity,
					Price: product.HETPrice,
					Total: float64(purchaseQuantity) * product.HETPrice,	
				}
				
				recomend = append(recomend, purchaseProduct)
			}
		}
	}


	res.Data = recomend
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}
