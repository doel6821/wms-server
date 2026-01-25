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

func (u *usecase) SaveCustomer(ctx context.Context, req cModels.RegisterCustomerRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	var customer pModels.Customers
	var err error
	if req.ID == 0 {
		customer, err = u.DB.GetPostgre().GetCustomerByPhone(ctx, req.Phone, req.Tenant)
		
		if err != nil && err.Error() == "record not found" {

			customer.Name = req.Name
			customer.Email = strings.ToLower(req.Email)
			customer.Phone = req.Phone
			customer.Address = req.Address
			customer.Tenant = strings.ToUpper(req.Tenant)
			customer.TermOfPayment = req.TermOfPayment
			customer.DiscountPercent = req.DiscountPercent
			customer.CancelOnBackOrder = req.CancelOnBackOrder

			err = u.DB.GetPostgre().SaveCustomers(ctx, customer)
			if err != nil {
				u.Logs.WithContext(ctx).WithError(err).Error("failed save customer")
				res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
				res.Meta.Message = "Gagal menyimpan data pelanggan"
				return res
			}
		} else if customer.ID != 0 {
			u.Logs.WithContext(ctx).WithError(err).Error("nomor telepon sudah terdaftar")
			res.Meta = helpers.GetMetaResponse(constants.RC_PHONE_NUMBER_ALREADY_USED)
			res.Meta.Message = "nomor telepon sudah terdaftar"
			return res
		}
	} else {
		customer = pModels.Customers{
			ID:                req.ID,
			Name:              req.Name,
			Email:             strings.ToLower(req.Email),
			Phone:             req.Phone,
			Address:           req.Address,
			Tenant:            strings.ToUpper(req.Tenant),
			TermOfPayment:     req.TermOfPayment,
			DiscountPercent:   req.DiscountPercent,
			CancelOnBackOrder: req.CancelOnBackOrder,
		}

		err = u.DB.GetPostgre().SaveCustomers(ctx, customer)
		if err != nil {
			u.Logs.WithContext(ctx).WithError(err).Error("failed save customer")
			res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
			res.Meta.Message = "Gagal menyimpan data pelanggan"
			return res
		}
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetCustomerList(ctx context.Context, tenant, name string, page, limit int) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listCustomer, total, err := u.DB.GetPostgre().GetCustomerList(ctx, tenant, name, page, limit)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get customer list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listCustomer
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetCustomerId(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	customer, err := u.DB.GetPostgre().GetCustomerById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get customer")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = customer
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) DeleteCustomerId(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	err := u.DB.GetPostgre().DeleteCustomerById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get customer")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}
