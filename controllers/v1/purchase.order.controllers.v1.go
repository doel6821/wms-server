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

// @Summary Save Purchase Order
// @Description Save Purchase Order
// @ID SavePurchaseOrder
// @Param body body cModels.PurchaseOrderRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Purchase-order [post]
func (c *v1Controller) SavePurchaseOrder(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.PurchaseOrderRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	tenant := ctx.GetString("tenant")
	res = c.Usecase.CreatePurchaseOrder(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save PurchaseOrder List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Purchase Order
// @Description Get List Purchase Order 
// @ID GetListPurchaseOrder
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Purchase-order/all [get]
func (c *v1Controller) ListPurchaseOrder(ctx *gin.Context) {
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

	res = c.Usecase.GetPurchaseOrderList(ctx, tenant, int64(supplierId), page, limit)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, supplierId, res, time.Since(trxTime))).Info("Get PurchaseOrder List")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get PurchaseOrder By ID
// @Description Get PurchaseOrder By ID
// @ID GetPurchaseOrderByID
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of PurchaseOrder"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Purchase-order/{id} [get]
func (c *v1Controller) GetPurchaseOrderById(ctx *gin.Context) {
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

	res = c.Usecase.GetPurchaseOrderDetailById(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Get PurchaseOrder By ID")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get Recomendation By Supplier ID
// @Description Get Recomendation By Supplier ID
// @ID GetRecomendationBySupplierID
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of Supplier"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/purchase-order/recomendation/{id} [get]
func (c *v1Controller) PurchaseOrderRecomendation(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var err error
	
	supplierId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return	
	}
	tenant := ctx.GetString("tenant")
	res = c.Usecase.GetPurchaseOrderRecomendation(ctx, tenant, int64(supplierId))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, supplierId, res, time.Since(trxTime))).Info("Get PurchaseOrder By ID")
	ctx.JSON(http.StatusOK, res)
}
