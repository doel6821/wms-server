package controllers

import (
	"net/http"

	"time"
	"wms-server/constants"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
	cModels "wms-server/controllers/v1/models"
	"github.com/gin-gonic/gin"
)


// @Summary Get Dashboard Data
// @Description Get Dashboard Data 
// @ID GetDashboardData
// @Param Authorization header string true "Bearer"
// @Param startDate query string true "StartDate"
// @Param endDate query string true "EndDate"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/dashboard [get]
func (c *v1Controller) GetDataDashboard(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.ReqListFinance
	
	if err := ctx.ShouldBindQuery(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	tenant := ctx.GetString("tenant")

	res = c.Usecase.GetDataDashboard(ctx, tenant, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, tenant, res, time.Since(trxTime))).Info("Get Account Payable List")
	ctx.JSON(http.StatusOK, res)
}

