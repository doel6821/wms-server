package controllers

import (
	"net/http"
	"strconv"
	"time"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"

	"github.com/gin-gonic/gin"
)

// @Summary Save Invoice
// @Description Save Invoice
// @ID SaveInvoice
// @Param body body cModels.InvoiceRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/invoice [post]
func (c *v1Controller) SaveInvoice(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.InvoiceRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	tenant := ctx.GetString("tenant")
	res = c.Usecase.CreateInvoice(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save Invoice List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Sales Order
// @Description Get List Sales Order 
// @ID GetListInvoice
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/invoice/all [get]
func (c *v1Controller) ListInvoice(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var page int
	var limit int
	var err error
	
	if ctx.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(ctx.Query("page"))
		if err != nil {
			res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
			ctx.JSON(http.StatusBadRequest, res)
			return	
		}
	}

	if ctx.Query("limit") == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
			ctx.JSON(http.StatusBadRequest, res)
			return	
		}
	}

	customerId, err := strconv.Atoi(ctx.Query("customerId"))
	if err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return	
	}
	tenant := ctx.GetString("tenant")

	res = c.Usecase.GetInvoiceList(ctx, tenant, int64(customerId), page, limit)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, customerId, res, time.Since(trxTime))).Info("Get Invoice List")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Invoice By ID
// @Description Get Invoice By ID
// @ID GetInvoiceByID
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of Invoice"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/invoice/{id} [get]
func (c *v1Controller) GetInvoiceDetailById(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var id int
	var err error
	
	id, err = strconv.Atoi(ctx.Param("id"))
	if err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return	
	}

	res = c.Usecase.GetInvoiceDetailById(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Get Invoice By ID")
	ctx.JSON(http.StatusOK, res)
}
