package usecases

import (
	"context"
	"strings"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
)

func (u *usecase) SaveSupplier(ctx context.Context, req cModels.RegisterSupplierRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	var supplier pModels.Supplier
	var err error
	if req.ID == 0 {
		supplier, err = u.DB.GetPostgre().GetSupplierByPhone(ctx, req.Phone, req.Tenant)
		if err != nil && err.Error() == "record not found" {

			supplier.Name = req.Name
			supplier.Email = strings.ToLower(req.Email)
			supplier.Phone = req.Phone
			supplier.Address = req.Address
			supplier.Tenant = strings.ToUpper(req.Tenant)
			supplier.TermOfPayment = req.TermOfPayment
			supplier.DiscountPercent = req.DiscountPercent

			err = u.DB.GetPostgre().SaveSupplier(ctx, supplier)
			if err != nil {
				u.Logs.WithContext(ctx).WithError(err).Error("failed save Supplier")
				res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
				return res
			}
		} else if supplier.ID != 0 {
			u.Logs.WithContext(ctx).WithError(err).Error("phone number already used")
			res.Meta = helpers.GetMetaResponse(constants.RC_PHONE_NUMBER_ALREADY_USED)
			return res
		}
	} else {
		supplier = pModels.Supplier{
			ID:              req.ID,
			Name:            req.Name,
			Email:           strings.ToLower(req.Email),
			Phone:           req.Phone,
			Address:         req.Address,
			Tenant:          strings.ToUpper(req.Tenant),
			TermOfPayment:   req.TermOfPayment,
			DiscountPercent: req.DiscountPercent,
		}

		err = u.DB.GetPostgre().SaveSupplier(ctx, supplier)
		if err != nil {
			u.Logs.WithContext(ctx).WithError(err).Error("failed save Supplier")
			res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
			return res
		}
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetSupplierList(ctx context.Context, tenant, name string, page, limit int) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listSupplier, total, err := u.DB.GetPostgre().GetSupplierList(ctx, tenant, name, page, limit)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Supplier list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listSupplier
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetSupplierId(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	Supplier, err := u.DB.GetPostgre().GetSupplierById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Supplier")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = Supplier
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) DeleteSupplierId(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	err := u.DB.GetPostgre().DeleteSupplierById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Supplier")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}
