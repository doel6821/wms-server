package usecases

import (
	"context"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
)

func (u *usecase) SaveProduct(ctx context.Context, tenant string, req cModels.RegisterProductRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	var product pModels.Product
	var err error
	if req.ID == 0 {
		product, err = u.DB.GetPostgre().GetProductByCode(ctx, tenant, req.Code)
		if err != nil && err.Error() == "record not found" {

			product.Name = req.Name
			product.Code = req.Code
			product.SupplierId = req.SupplierId
			product.HETPrice = req.HETPrice
			product.CostPrice = req.CostPrice
			product.AvgPrice = req.AvgPrice
			product.StockAllocation = req.StockAllocation
			product.StockOnHand = req.StockOnHand
			product.StockOnPurchase = req.StockOnPurchase
			product.StockOnReceive = req.StockOnReceive
			product.StockPacking = req.StockPacking
			product.LeadTimeDays = req.LeadTimeDays
			product.Tenant = tenant

			err = u.DB.GetPostgre().SaveProduct(ctx, product)
			if err != nil {
				u.Logs.WithContext(ctx).WithError(err).Error("failed save Product")
				res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
				return res
			}
		} else if product.ID != 0 {
			u.Logs.WithContext(ctx).WithError(err).Error("phone number already used")
			res.Meta = helpers.GetMetaResponse(constants.RC_PHONE_NUMBER_ALREADY_USED)
			return res
		}
	} else {
		product = pModels.Product{
			ID:              req.ID,
			Name:            req.Name,
			Code:            req.Code,
			SupplierId:      req.SupplierId,
			HETPrice:        req.HETPrice,
			CostPrice:       req.CostPrice,
			AvgPrice:        req.AvgPrice,
			StockAllocation: req.StockAllocation,
			StockOnHand:     req.StockOnHand,
			StockOnPurchase: req.StockOnPurchase,
			StockOnReceive:  req.StockOnReceive,
			StockPacking:    req.StockPacking,
			LeadTimeDays:    req.LeadTimeDays,
			Tenant:          tenant,
		}

		err = u.DB.GetPostgre().SaveProduct(ctx, product)
		if err != nil {
			u.Logs.WithContext(ctx).WithError(err).Error("failed save Product")
			res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
			return res
		}
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetProductList(ctx context.Context, tenant, name, code, supplier string, page, limit int) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listProduct, total, err := u.DB.GetPostgre().GetProductList(ctx, name, tenant, code, supplier, page, limit)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Product list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listProduct
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetProductId(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	Product, err := u.DB.GetPostgre().GetProductById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Product")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = Product
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) DeleteProductId(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	err := u.DB.GetPostgre().DeleteProductById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Product")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}
