package controllers

import (
	"net/http"
	"strconv"
	"time"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"

	"github.com/gin-gonic/gin"
)


// @Summary Get Available Receive Order
// @Description Get Available Receive Order
// @ID GetAvaliableReceiveOrder
// @Param supplierId query string false "SupplierId"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/receive-order [get]
func (c *v1Controller) GetAvailableReceive(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var err error

	supplierIdStr := ctx.Query("supplierId")
	productIdStr := ctx.Query("productId")
	var supplierId int
	if supplierIdStr != "" {
		supplierId, err = strconv.Atoi(supplierIdStr)
		if err != nil {
			res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
			ctx.JSON(http.StatusBadRequest, res)
			return	
		}
	} 
	var productId int
	if productIdStr != "" {
		productId, err = strconv.Atoi(productIdStr)
		if err != nil {
			res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
			ctx.JSON(http.StatusBadRequest, res)
			return	
		}
	} 

	tenant := ctx.GetString("tenant")

	res = c.Usecase.GetAvailableReceiveList(ctx, tenant, int64(supplierId), int64(productId))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, supplierId, res, time.Since(trxTime))).Info("Get ReceiveOrder List")
	ctx.JSON(http.StatusOK, res)

}


// @Summary Save Receive Order
// @Description Save Receive Order
// @ID SaveReceiveOrder
// @Param body body cModels.ReceiveOrderRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/receive-order/ [post]
func (c *v1Controller) SaveReceiveOrder(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.ReceiveOrderRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	tenant := ctx.GetString("tenant")
	res = c.Usecase.CreateReceive(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save ReceiveOrder List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Receive Order
// @Description Get List Receive Order 
// @ID GetListReceiveOrder
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/receive-order/all [get]
func (c *v1Controller) ListReceiveOrder(ctx *gin.Context) {
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

	supplierIdStr := ctx.Query("supplierId")
	dueDate := ctx.Query("dueDate")
	var supplierId int
	if supplierIdStr != "" {
		supplierId, err = strconv.Atoi(supplierIdStr)
		if err != nil {
			res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
			ctx.JSON(http.StatusBadRequest, res)
			return	
		}
	}
	tenant := ctx.GetString("tenant")

	res = c.Usecase.GetReceiveList(ctx, tenant, int64(supplierId), page, limit, dueDate)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, supplierId, res, time.Since(trxTime))).Info("Get ReceiveOrder List")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get ReceiveOrder By ID
// @Description Get ReceiveOrder By ID
// @ID GetReceiveOrderByID
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of ReceiveOrder"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Receive-order/{id} [get]
func (c *v1Controller) GetReceiveOrderDetailById(ctx *gin.Context) {
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

	res = c.Usecase.GetReceiveDetailById(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Get ReceiveOrder By ID")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get ReceiveOrder By ProductID
// @Description Get ReceiveOrder By ProductID
// @ID GetReceiveOrderByID
// @Param Authorization header string true "Bearer"
// @Param productId path int true "ProductID of ReceiveOrder"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Receive-order/{productId} [get]
func (c *v1Controller) GetReceiveOrderDetailByProductId(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var id int
	var err error
	
	id, err = strconv.Atoi(ctx.Param("productId"))
	if err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return	
	}

	res = c.Usecase.GetReceiveDetailByProductId(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Get ReceiveOrder By ID")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Save Stocked
// @Description Save Stocked
// @ID SaveStocked
// @Param body body cModels.StockedRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/receive-order/stocked [post]
func (c *v1Controller) SaveStocked(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.StockedRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	tenant := ctx.GetString("tenant")
	res = c.Usecase.CreateStocked(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save Stocked")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Update Receive Order
// @Description Update Receive Order
// @ID UpdateReceiveOrder
// @Param body body pModels.ReceiveOrder true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/receive-order/ [put]
func (c *v1Controller) UpdateReceiveOrder(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req pModels.ReceiveOrder


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	tenant := ctx.GetString("tenant")
	res = c.Usecase.UpdateReceive(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save ReceiveOrder List")
	ctx.JSON(http.StatusOK, res)
}



