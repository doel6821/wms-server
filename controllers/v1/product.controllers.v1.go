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

// @Summary Save Product
// @Description Save Product
// @ID SaveProduct
// @Param body body cModels.RegisterProductRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/product [post]
func (c *v1Controller) SaveProduct(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.RegisterProductRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res = c.Usecase.SaveProduct(ctx, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save Product List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Product
// @Description Get List Product 
// @ID GetListProduct
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/product/all [get]
func (c *v1Controller) ListProduct(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var page int
	var limit int
	var name string
	var code string
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
	code = ctx.Query("code")
	tenant := ctx.GetString("tenant")
	
	res = c.Usecase.GetProductList(ctx, tenant, name, code, page, limit)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, name, res, time.Since(trxTime))).Info("Get Product List")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Product By ID
// @Description Get Product By ID
// @ID GetProductByID
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of Product"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/product/{id} [get]
func (c *v1Controller) GetProductById(ctx *gin.Context) {
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

	res = c.Usecase.GetProductId(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Get Product By ID")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Product
// @Description Delete Product 
// @ID DeleteProduct
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of Product"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/product/{id} [delete]
func (c *v1Controller) DeleteProductById(ctx *gin.Context) {
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

	res = c.Usecase.DeleteProductId(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Delete Product")
	ctx.JSON(http.StatusOK, res)
}