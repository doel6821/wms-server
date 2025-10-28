package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"

	"github.com/gin-gonic/gin"
)

// @Summary Save Location
// @Description Save Location
// @ID SaveLocation
// @Param body body cModels.RegisterLocationRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/location [post]
func (c *v1Controller) SaveLocation(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.RegisterLocationRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	tenant := ctx.GetString("tenant")
	res = c.Usecase.SaveLocation(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save Location List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Location
// @Description Get List Location 
// @ID GetListLocation
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Location/all [get]
func (c *v1Controller) ListLocation(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var name string
	
	tenant := ctx.GetString("tenant")

	var query cModels.LocationQueryParams

	// Bind query string ke struct
	if err := ctx.ShouldBindQuery(&query); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	fmt.Println(query)
	
	res = c.Usecase.GetLocationList(ctx, tenant, query)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, name, res, time.Since(trxTime))).Info("Get Location List")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Location By ID
// @Description Get Location By ID
// @ID GetLocationByID
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of Location"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Location/{id} [get]
func (c *v1Controller) GetLocationById(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	
	id := ctx.Param("id")
	tenant := ctx.GetString("tenant")
	idNumber, err := strconv.Atoi(id)
	if err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return	
	}
	res = c.Usecase.GetLocationById(ctx,tenant, int64(idNumber))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, tenant, res, time.Since(trxTime))).Info("Get Location By code")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get Location By Code
// @Description Get Location By Code
// @ID GetLocationByCode
// @Param Authorization header string true "Bearer"
// @Param code path int true "Code of Location"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Location/{code} [get]
func (c *v1Controller) GetLocationByCode(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	
	code := ctx.Param("code")
	tenant := ctx.GetString("tenant")

	res = c.Usecase.GetLocationByCode(ctx,tenant, code)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, tenant, res, time.Since(trxTime))).Info("Get Location By code")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Location
// @Description Delete Location 
// @ID DeleteLocation
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of Location"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Location/{id} [delete]
func (c *v1Controller) DeleteLocationById(ctx *gin.Context) {
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

	res = c.Usecase.DeleteLocationId(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Delete Location")
	ctx.JSON(http.StatusOK, res)
}