package usecases

import (
	"context"
	"strings"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
	uModels "wms-server/usecases/v1/models"
	pModels "wms-server/databases/postgre/models"
)


func (u *usecase) Login(ctx context.Context, req cModels.LoginRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR),
	}

	user, err := u.DB.GetPostgre().FindUser(ctx, req.Email)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("Error get existing data")
		return res
	}

	if user.Password != req.Password {
		u.Logs.WithContext(ctx).WithError(err).Error("invalid user name or password")
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_INVALID_EMAIL_OR_PASSWORD)
		return res
	} 

	token, expired, err := helpers.GenerateToken(user)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("Error generate token")
		return res
	}
	res.Data = uModels.AuthResponse{
		Token: token,
		TokenType: "Bearer",
		ExpiresIn: expired,
	}
	res.Meta = helpers.GetNewMetaResponse("en", constants.RC_SUCCESS)
	return res
}

func (u *usecase) RegisterUser(ctx context.Context, req cModels.RegisterRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR),
	}

	user, err := u.DB.GetPostgre().FindUser(ctx, strings.ToLower(req.Email))
	if err != nil && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("Error get existing data")
		return res
	}

	if user.ID != 0 {
		u.Logs.WithContext(ctx).WithError(err).Error("email already register")
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_EMAIL_ALREADY_USED)
		return res
	}

	user = pModels.User {
		Email: strings.ToLower(req.Email),
		Role: req.Role,
		Tenant: strings.ToUpper(req.Tenant),
		Password: req.Password,
	}

	err = u.DB.GetPostgre().Save(ctx, user)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed save user")
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR)
		return res
	}
	res.Meta = helpers.GetNewMetaResponse("en", constants.RC_SUCCESS)
	return res	
}


func (u *usecase) RegisterTenant(ctx context.Context, req cModels.RegisterRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR),
	}

	user, err := u.DB.GetPostgre().FindUser(ctx, strings.ToLower(req.Email))
	if err != nil  && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("Error get existing data")
		return res
	}

	tenant, err := u.DB.GetPostgre().FindTenant(ctx, strings.ToUpper(req.Tenant))
	if err != nil  && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("Error get existing data")
		return res
	}

	if tenant.ID != 0 {
		u.Logs.WithContext(ctx).WithError(err).Error("store name already register")
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_STORE_NAME_ALREADY_USED)
		return res
	}

	if user.ID != 0 {
		u.Logs.WithContext(ctx).WithError(err).Error("email already register")
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_EMAIL_ALREADY_USED)
		return res
	}
	

	user = pModels.User {
		Email: strings.ToLower(req.Email),
		Role: req.Role,
		Tenant: strings.ToUpper(req.Tenant),
		Password: req.Password,
	}

	err = u.DB.GetPostgre().Save(ctx, user)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed save user")
		res.Meta = helpers.GetNewMetaResponse("en", constants.RC_GENERAL_ERROR)
		return res
	}
	res.Meta = helpers.GetNewMetaResponse("en", constants.RC_SUCCESS)
	return res	
}