package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	authorizationSecretKey = helpers.GetEnv("SECRET_KEY")
	secretKey              = []byte(authorizationSecretKey)
)

// @Summary Authenticate - GenerateToken
// @Description Authenticate 
// @ID Authenticate
// @Param body body cModels.LoginRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/login [post]
func (c *v1Controller) Login(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.LoginRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res = c.Usecase.Login(ctx, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Login")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Authenticate - Register
// @Description Register New User 
// @ID Register
// @Param body body cModels.RegisterRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1.0.0/register-user [post]
func (c *v1Controller) RegisterUser(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.RegisterRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res = c.Usecase.RegisterUser(ctx, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Register User")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Authenticate - Register
// @Description Register New User 
// @ID Register
// @Param body body cModels.RegisterRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1.0.0/register-tenant [post]
func (c *v1Controller) RegisterTenant(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.RegisterRequest


	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res = c.Usecase.RegisterTenant(ctx, req)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Register Tenant")
	ctx.JSON(http.StatusOK, res)
}


func (c *v1Controller) CekToken(ctx *gin.Context) {
	var res hModels.Response
	authHeader := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		res.Meta = helpers.GetMetaResponse(constants.RC_UNAUTHORIZED)
		ctx.JSON(http.StatusOK, res)
		return
	}
	newTokenString := strings.Split(authHeader, " ")
	newData := &cModels.ClaimsUserDto{}
	_, err := jwt.ParseWithClaims(newTokenString[1], newData, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unexpected signing method"})
			ctx.Abort()
		}

		return secretKey, nil
	})
	fmt.Println("err :" ,err)
	v, _ := err.(*jwt.ValidationError)

	if v != nil {
		if v.Errors == jwt.ValidationErrorExpired && newData.ExpiresAt > time.Now().Unix() {
			res.Meta = helpers.GetMetaResponse(constants.RC_EXPIRED_TOKEN)
			ctx.JSON(http.StatusUnauthorized, res)
			ctx.Abort()
			return
		}
	}

	if err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_INVALID_TOKEN)
		ctx.JSON(http.StatusUnauthorized, res)
		ctx.Abort()
		return
	}

	ctx.Set("email", newData.Scope.Email)
	ctx.Set("role", newData.Scope.Role)
	ctx.Set("tenant", newData.Scope.Tenant)
	
	ctx.Next()
}
