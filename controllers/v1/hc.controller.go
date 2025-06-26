package controllers

import (
	"wms-server/controllers/v1/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"

	"github.com/gin-gonic/gin"
)

// HealtCheck godoc
// @Summary HealtCheck
// @Description hc endpoint
// @Tags HC
// @ID hc
// @Produce json
// @Success 200 {object} hModels.Response
// @Failure 400 {object} hModels.Response
// @Failure 401 {object} hModels.Response
// @Router /healthcheck [get]
// HealthCheck ...
func (c *v1Controller) HealthCheck(ctx *gin.Context) {
	var res hModels.Response

	head := models.Header{}

	ctx.ShouldBindHeader(&head)

	res = c.Usecase.HealthCheck(ctx)

	ctx.JSON(helpers.MapStatusCode(res.Meta), res)
}
