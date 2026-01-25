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
// @ID SaveSalesOrder
// @Param body body cModels.SalesOrderRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/sales-order [post]
func (c *v1Controller) SaveSalesOrder(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.SalesOrderRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	tenant := ctx.GetString("tenant")
	res = c.Usecase.CreateSalesOrder(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save SalesOrder List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Sales Order
// @Description Get List Sales Order 
// @ID GetListSalesOrder
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/sales-order/all [get]
func (c *v1Controller) ListSalesOrder(ctx *gin.Context) {
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

	customerIdStr := ctx.Query("customerId")
	allocation := ctx.Query("allocation")
	var customerId int
	if customerIdStr != "" {
		customerId, err = strconv.Atoi(customerIdStr)
		if err != nil {
			res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
			ctx.JSON(http.StatusBadRequest, res)
			return	
		}
	}
	tenant := ctx.GetString("tenant")

	res = c.Usecase.GetSalesOrderList(ctx, tenant, allocation, int64(customerId), page, limit)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, customerId, res, time.Since(trxTime))).Info("Get SalesOrder List")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get SalesOrder By ID
// @Description Get SalesOrder By ID
// @ID GetSalesOrderByID
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of SalesOrder"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/sales-order/{id} [get]
func (c *v1Controller) GetSalesOrderById(ctx *gin.Context) {
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

	res = c.Usecase.GetSalesOrderDetailById(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Get SalesOrder By ID")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get SalesOrderItem By ID
// @Description Get SalesOrderItem By ID
// @ID GetSalesOrderItemByID
// @Param Authorization header string true "Bearer"
// @Param productId query int true "productId"
// @Param type query int true "type"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/sales-order-item [get]
func (c *v1Controller) GetSalesOrderItemByProductId(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var err error

	productIdStr := ctx.Query("productId")
	tipe := ctx.Query("tipe")
	
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return	
	}

	res = c.Usecase.GetSalesOrderItemByProductId(ctx, int64(productId), tipe)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, productId, res, time.Since(trxTime))).Info("Get SalesOrder By ID")
	ctx.JSON(http.StatusOK, res)
}
