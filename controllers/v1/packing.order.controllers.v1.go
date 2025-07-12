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

// @Summary Save Sales Order
// @Description Save Sales Order
// @ID SavePackingOrder
// @Param body body cModels.PackingOrderRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/packing-order [post]
func (c *v1Controller) SavePackingOrder(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.PackingOrderRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	tenant := ctx.GetString("tenant")
	res = c.Usecase.CreatePackingOrder(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save PackingOrder List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Sales Order
// @Description Get List Sales Order 
// @ID GetListPackingOrder
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/packing-order/all [get]
func (c *v1Controller) ListPackingOrder(ctx *gin.Context) {
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

	res = c.Usecase.GetPackingOrderList(ctx, tenant, int64(customerId), page, limit)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, customerId, res, time.Since(trxTime))).Info("Get PackingOrder List")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get PackingOrder By ID
// @Description Get PackingOrder By ID
// @ID GetPackingOrderByID
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of PackingOrder"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/packing-order/{id} [get]
func (c *v1Controller) GetPackingOrderDetailById(ctx *gin.Context) {
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

	res = c.Usecase.GetPackingOrderDetailById(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Get PackingOrder By ID")
	ctx.JSON(http.StatusOK, res)
}
