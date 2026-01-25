package controllers

import (
	"net/http"

	"time"
	"wms-server/constants"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
	cModels "wms-server/controllers/v1/models"
	"github.com/gin-gonic/gin"
)

// @Summary Save AccountPayable
// @Description Save AccountPayable
// @ID SaveAccountPayable
// @Param body body cModels.AccountPayableRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/finance/payment [post]
func (c *v1Controller) SaveAccountPayable(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.ReqAccountPayable

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	tenant := ctx.GetString("tenant")
	res = c.Usecase.CreateAccountPayable(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save Invoice List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Account Payable
// @Description Get List Account Payable 
// @ID GetListAccounPayable
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/finance/payment/all [get]
func (c *v1Controller) ListAccountPayable(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.ReqListFinance
	
	if err := ctx.ShouldBindQuery(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	tenant := ctx.GetString("tenant")

	res = c.Usecase.GetAccountPayableList(ctx, tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, tenant, res, time.Since(trxTime))).Info("Get Account Payable List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Save AccountReceivable
// @Description Save AccountReceivable
// @ID SaveAccountReceivable
// @Param body body cModels.AccountReceivableRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/finance/receive [post]
func (c *v1Controller) SaveAccountReceivable(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.ReqAccountReceiveble

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	tenant := ctx.GetString("tenant")
	res = c.Usecase.CreateAccountReceivable(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save Invoice List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Account Receivable
// @Description Get List Account Receivable 
// @ID GetListAccounReceivable
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/finance/receive/all [get]
func (c *v1Controller) ListAccountReceivable(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.ReqListFinance
	
	if err := ctx.ShouldBindQuery(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	tenant := ctx.GetString("tenant")

	res = c.Usecase.GetAccountReceivableList(ctx, tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, tenant, res, time.Since(trxTime))).Info("Get Account Payable List")
	ctx.JSON(http.StatusOK, res)
}

