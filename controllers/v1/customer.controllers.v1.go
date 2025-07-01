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

// @Summary Save Customer
// @Description Save Customer
// @ID SaveCustomer
// @Param body body cModels.RegisterCustomerRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/customer [post]
func (c *v1Controller) SaveCustomer(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.RegisterCustomerRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res = c.Usecase.SaveCustomer(ctx, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save Customer List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Customer
// @Description Get List Customer 
// @ID GetListCustomer
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/customer/all [get]
func (c *v1Controller) ListCustomer(ctx *gin.Context) {
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

	res = c.Usecase.GetCustomerList(ctx, name, page, limit)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, name, res, time.Since(trxTime))).Info("Get Customer List")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get List Customer
// @Description Get List Customer 
// @ID GetListCustomer
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of customer"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/customer/{id} [get]
func (c *v1Controller) GetCustomerById(ctx *gin.Context) {
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

	res = c.Usecase.GetCustomerId(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Get Customer By ID")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Customer
// @Description Delete Customer 
// @ID DeleteCustomer
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of customer"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/customer/{id} [delete]
func (c *v1Controller) DeleteCustomerById(ctx *gin.Context) {
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

	res = c.Usecase.DeleteCustomerId(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Delete Customer")
	ctx.JSON(http.StatusOK, res)
}