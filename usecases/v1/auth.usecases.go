package usecases

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
	uModels "wms-server/usecases/v1/models"
)

func (u *usecase) Login(ctx context.Context, req cModels.LoginRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	user, err := u.DB.GetPostgre().FindUser(ctx, req.Email)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("Error get existing data")
		return res
	}

	if user.Password != req.Password {
		u.Logs.WithContext(ctx).WithError(err).Error("invalid user name or password")
		res.Meta = helpers.GetMetaResponse(constants.RC_INVALID_EMAIL_OR_PASSWORD)
		return res
	}

	token, expired, err := helpers.GenerateToken(user)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("Error generate token")
		return res
	}
	res.Data = uModels.AuthResponse{
		Token:     token,
		Role:      user.Role,
		TokenType: "Bearer",
		ExpiresIn: expired,
	}
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) RegisterUser(ctx context.Context, req cModels.RegisterRequest, tenant string) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	user, err := u.DB.GetPostgre().FindUser(ctx, strings.ToLower(req.Email))
	if err != nil && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("Error get existing data")
		return res
	}

	if user.ID != 0 {
		u.Logs.WithContext(ctx).WithError(err).Error("email already register")
		res.Meta = helpers.GetMetaResponse(constants.RC_EMAIL_ALREADY_USED)
		return res
	}

	user = pModels.User{
		Email:    strings.ToLower(req.Email),
		Role:     req.Role,
		Tenant:   tenant,
		Password: req.Password,
	}

	err = u.DB.GetPostgre().Save(ctx, user)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed save user")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) RegisterTenant(ctx context.Context, req cModels.RegisterRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	user, err := u.DB.GetPostgre().FindUser(ctx, strings.ToLower(req.Email))
	if err != nil && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("Error get existing data")
		return res
	}

	tenant, err := u.DB.GetPostgre().FindTenant(ctx, strings.ToUpper(req.Tenant))
	if err != nil && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("Error get existing data")
		return res
	}

	if tenant.ID != 0 {
		u.Logs.WithContext(ctx).WithError(err).Error("store name already register")
		res.Meta = helpers.GetMetaResponse(constants.RC_STORE_NAME_ALREADY_USED)
		return res
	}

	if user.ID != 0 {
		u.Logs.WithContext(ctx).WithError(err).Error("email already register")
		res.Meta = helpers.GetMetaResponse(constants.RC_EMAIL_ALREADY_USED)
		return res
	}

	user = pModels.User{
		Email:    strings.ToLower(req.Email),
		Role:     "admin",
		Tenant:   strings.ToUpper(req.Tenant),
		Password: req.Password,
	}

	err = u.DB.GetPostgre().Save(ctx, user)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed save user")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) UpdateUser(ctx context.Context, req cModels.RegisterRequest, tenant string, id int) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	user, err := u.DB.GetPostgre().FindUserById(ctx, id)
	if err != nil && err.Error() != "record not found" {
		u.Logs.WithContext(ctx).WithError(err).Error("user not register")
		return res
	}

	user.Role = req.Role

	err = u.DB.GetPostgre().Save(ctx, user)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed save user")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetUserList(ctx context.Context, tenant, email string, page, limit int) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listUser, total, err := u.DB.GetPostgre().GetUserList(ctx, tenant, email, page, limit)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get User list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listUser
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}


func (u *usecase) ForgotPassword(ctx context.Context, req cModels.LoginRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	user, err := u.DB.GetPostgre().FindUser(ctx, req.Email)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("Error get existing data")
		return res
	}

	random := helpers.RandomString(10)
	sum := sha512.Sum512([]byte(random))
	newPassword := hex.EncodeToString(sum[:])
	fmt.Println("newPassword : ", random)
	fmt.Println(newPassword)
	user.Password = newPassword

	err = u.DB.GetPostgre().Save(ctx, user)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed save user")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	to := req.Email
	subject := "Reset Password"
	templateEmail := `
		<p>Hi, Kami telah menerima permintaanmu untuk melakukan reset password akun WMS.</p>
		<p>Silahkan login ulang menggunakan username dan password dibawah: </p>
		<p>Username: ` + req.Email + `</p>
		<p>Password: ` + newPassword + `</p>
		<br>
		<p>Abaikan email ini jika kamu tidak pernah meminta untuk melakukan reset password.</p>`

	err = u.MailSmpt.SendEmailSMTP(to, "" ,subject, templateEmail)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed send email")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}


func (u *usecase) ChangePassword(ctx context.Context, req cModels.ChangePasswordRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	user, err := u.DB.GetPostgre().FindUser(ctx, req.Email)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("Error get existing data")
		return res
	}

	if user.Password != req.OldPassword {
		u.Logs.WithContext(ctx).WithError(err).Error("password tidak sesuai")
		res.Meta = helpers.GetMetaResponse(constants.RC_BADREQUEST)
		return res
	}

	
	user.Password = req.NewPassword

	err = u.DB.GetPostgre().Save(ctx, user)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed save user")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res

}
