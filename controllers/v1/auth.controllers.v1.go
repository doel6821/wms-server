package controllers

import (
	"fmt"
	"net/http"
	"strconv"
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
	req.Tenant = ctx.GetString("tenant")
	res = c.Usecase.RegisterUser(ctx, req, req.Tenant)
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


// @Summary Authenticate - Update
// @Description Update User 
// @ID Update
// @Param body body cModels.RegisterRequest true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1.0.0/update-user [post]
func (c *v1Controller) UpdateUser(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var req cModels.RegisterRequest

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return	
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	req.Tenant = ctx.GetString("tenant")
	res = c.Usecase.UpdateUser(ctx, req, req.Tenant, id)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, req, res, time.Since(trxTime))).Info("Update User")
	ctx.JSON(http.StatusOK, res)
}


// @Summary Get List User
// @Description Get List User 
// @ID GetListUser
// @Param Authorization header string true "Bearer"
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Param name query string false "Name"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /wms-server/v.1/user/all [get]
func (c *v1Controller) ListUser(ctx *gin.Context) {
	trxTime := time.Now()
	var res hModels.Response
	var page int
	var limit int
	var email string
	var err error
	
	if ctx.Query("page") == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(ctx.Query("page"))
		if err != nil {
			res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
			fmt.Println("apakah dr sini 1")
			ctx.JSON(http.StatusBadRequest, res)
			return	
		}
	}

	if ctx.Query("limit") == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			fmt.Println("apakah dr sini 2")
			res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
			ctx.JSON(http.StatusBadRequest, res)
			return	
		}
	}

	email = ctx.Query("email")
	tenant := ctx.GetString("tenant")
	fmt.Println(tenant, "=======>>>>")
	res = c.Usecase.GetUserList(ctx, tenant, email, page, limit)
	c.Logs.WithFields(helpers.GettingResponseLog( ctx, email, res, time.Since(trxTime))).Info("Get Customer List")
	ctx.JSON(http.StatusOK, res)
}