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

// @Summary Save Configuration
// @Description Save Configuration
// @ID SaveConfiguration
// @Param body body cModels.RegisterConfigurationRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/Configuration [post]
func (c *v1Controller) SaveConfiguration(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.RegisterConfigurationRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	tenant := ctx.GetString("tenant")
	res = c.Usecase.SaveConfiguration(ctx,tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Save Configuration List")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List Configuration
// @Description Get List Configuration 
// @ID GetListConfiguration
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/config/all [get]
func (c *v1Controller) ListConfiguration(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var name string
	
	tenant := ctx.GetString("tenant")

	var query cModels.ConfigurationQueryParams

	// Bind query string ke struct
	if err := ctx.ShouldBindQuery(&query); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	
	res = c.Usecase.GetConfigurationList(ctx, tenant, query)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, name, res, time.Since(trxTime))).Info("Get Configuration List")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Configuration By ID
// @Description Get Configuration By ID
// @ID GetConfigurationByID
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of Configuration"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/config/{id} [get]
func (c *v1Controller) GetConfigurationById(ctx *gin.Context) {
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
	res = c.Usecase.GetConfigurationById(ctx,tenant, int64(idNumber))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, tenant, res, time.Since(trxTime))).Info("Get Configuration By code")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get Configuration By Name
// @Description Get Configuration By Name
// @ID GetConfigurationByName
// @Param Authorization header string true "Bearer"
// @Param name path int true "Name of Configuration"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/config/{name} [get]
func (c *v1Controller) GetConfigurationByName(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	
	name := ctx.Param("name")
	tenant := ctx.GetString("tenant")

	res = c.Usecase.GetConfigurationByName(ctx,tenant, name)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, tenant, res, time.Since(trxTime))).Info("Get Configuration By code")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Configuration
// @Description Delete Configuration 
// @ID DeleteConfiguration
// @Param Authorization header string true "Bearer"
// @Param id path int true "ID of Configuration"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/config/{id} [delete]
func (c *v1Controller) DeleteConfigurationById(ctx *gin.Context) {
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

	res = c.Usecase.DeleteConfigurationId(ctx, int64(id))
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, id, res, time.Since(trxTime))).Info("Delete Configuration")
	ctx.JSON(http.StatusOK, res)
}