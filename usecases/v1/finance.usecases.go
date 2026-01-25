package usecases

import (
	"context"
	"time"
	"wms-server/constants"
	pModels "wms-server/databases/postgre/models"
	cModels "wms-server/controllers/v1/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
)


func (u *usecase)CreateAccountPayable(ctx context.Context, tenant string, req cModels.ReqAccountPayable) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	ap := pModels.AccountPayable{
		ReceiveId       : req.ReceiveId,
		Amount          : req.Amount,
		PaymentDate     : time.Now(),
		PaymentMethod   : req.PaymentMethod,
		ReferenceNumber : req.ReferenceNumber,
		Notes           : req.Notes,
		Tenant          : tenant,
		InvoiceNumber   : req.InvoiceNumber,
	}

	err := u.DB.GetPostgre().SaveAccountPayable(ctx, ap)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed save account payable")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	err = u.DB.GetPostgre().UpdatePaymentStatusReceiveOrder(ctx, req.ReceiveId, "Lunas", time.Now() )
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed updatePaymentStatusReceiveOrder")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res

}

func (u *usecase)CreateAccountReceivable(ctx context.Context, tenant string, req cModels.ReqAccountReceiveble) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	ar := pModels.AccountReceiveble{
		InvoiceId       : req.InvoiceId,
		Amount          : req.Amount,
		PaymentDate     : time.Now(),
		PaymentMethod   : req.PaymentMethod,
		ReferenceNumber : req.ReferenceNumber,
		Notes           : req.Notes,
		Tenant          : tenant,
		InvoiceNumber   : req.InvoiceNumber,
	}

	err := u.DB.GetPostgre().SaveAccountReceivable(ctx, ar)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed save customer")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	err = u.DB.GetPostgre().UpdatePaymentStatusInvoiceOrder(ctx, req.InvoiceId, "Lunas", time.Now() )
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed updatePaymentStatusReceiveOrder")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res

}

func (u *usecase)GetAccountPayableList(ctx context.Context,tenant string, req cModels.ReqListFinance) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listAP, total, err := u.DB.GetPostgre().GetAccountPayableList(ctx, tenant, req)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get account payable list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Data = listAP
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res

}

func (u *usecase)GetAccountReceivableList(ctx context.Context,tenant string, req cModels.ReqListFinance) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listAR, total, err := u.DB.GetPostgre().GetAccountReceivableList(ctx, tenant, req)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get customer list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Data = listAR
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res


}