package usecases

import (
	"context"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
)


func (u *usecase)GetDataDashboard(ctx context.Context,tenant string, req cModels.ReqListFinance) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	startDate := req.StartDate + " 00:00:00"
	endDate := req.EndDate + " 23:59:59"

	sales, err := u.DB.GetPostgre().GetSalesOrderTotal(ctx, tenant, startDate, endDate)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get sales order total")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	purchase, err := u.DB.GetPostgre().GetPurchaseOrderTotal(ctx, tenant, startDate , endDate)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get purchase order total")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	ar, err := u.DB.GetPostgre().GetAccountReceivableTotal(ctx, tenant, req)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get account receaivable total")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	ap, err := u.DB.GetPostgre().GetAccountPayableTotal(ctx, tenant, req)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get account payable total")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	resData := map[string]interface{}{}

	resData["sales"] = sales
	resData["purchase"] = purchase
	resData["totalReceivePaymentCustomer"] = ar
	resData["totalPaymentSupplier"] = ap

	res.Data = resData
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res

}

