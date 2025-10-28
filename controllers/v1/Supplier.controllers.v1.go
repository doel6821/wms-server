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

// @Summary Save Supplier
// @Description Save Supplier
// @ID SaveSupplier
// @Param body body cModels.RegisterSupplierRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Supplier [post]
func (c *v1Controller) SaveSupplier(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.RegisterSupplierRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	req.Tenant = ctx.GetString("tenant")
	res = c.Usecase.SaveSupplier(ctx, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save Supplier List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Supplier
// @Description Get List Supplier 
// @ID GetListSupplier
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/supplier/all [get]
func (c *v1Controller) ListSupplier(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var page int
	var limit int
	var name string
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

	name = ctx.Query("name")
	tenant := ctx.GetString("tenant")
	res = c.Usecase.GetSupplierList(ctx, tenant, name, page, limit)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, name, res, time.Since(trxTime))).Info("Get Supplier List")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get List Supplier
// @Description Get List Supplier 
// @ID GetListSupplier
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of Supplier"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Supplier/{id} [get]
func (c *v1Controller) GetSupplierById(ctx *gin.Context) {
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

	res = c.Usecase.GetSupplierId(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Get Supplier By ID")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Supplier
// @Description Delete Supplier 
// @ID DeleteSupplier
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of Supplier"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Supplier/{id} [delete]
func (c *v1Controller) DeleteSupplierById(ctx *gin.Context) {
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

	res = c.Usecase.DeleteSupplierId(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Delete Supplier")
	ctx.JSON(http.StatusOK, res)
}