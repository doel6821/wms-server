package usecases

import (
	"context"
	"strings"
	"wms-server/constants"
	cModels "wms-server/controllers/v1/models"
	pModels "wms-server/databases/postgre/models"
	"wms-server/helpers"
	hModels "wms-server/helpers/models"
)

func (u *usecase) SaveLocation(ctx context.Context, tenant string, req cModels.RegisterLocationRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	var location pModels.Location
	var err error
	if req.ID == 0 {
		location, err = u.DB.GetPostgre().GetLocationByLocationCode(ctx, tenant, req.Code)
		if err != nil && err.Error() == "record not found" {

			location.Name = req.Name
			location.Code = req.Code
			location.Tenant = strings.ToUpper(tenant)

			err = u.DB.GetPostgre().SaveLocation(ctx, location)
			if err != nil {
				u.Logs.WithContext(ctx).WithError(err).Error("failed save Location")
				res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
				res.Meta.Message = "Gagal menyimpan data lokasi"
				return res
			}
		} else if location.ID != 0 {
			u.Logs.WithContext(ctx).WithError(err).Error("location code already exist")
			res.Meta = helpers.GetMetaResponse(constants.RC_LOCATION_EXIST)
			res.Meta.Message = "kode lokasi sudah digunakan"
			return res
		}
	} else {
		location = pModels.Location{
			ID:     req.ID,
			Name:   req.Name,
			Code:   req.Code,
			Tenant: strings.ToUpper(tenant),
		}

		err = u.DB.GetPostgre().SaveLocation(ctx, location)
		if err != nil {
			u.Logs.WithContext(ctx).WithError(err).Error("failed save Location")
			res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
			res.Meta.Message = "Gagal menyimpan data lokasi"
			return res
		}
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetLocationList(ctx context.Context, tenant string, query cModels.LocationQueryParams) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listLocation, total, err := u.DB.GetPostgre().GetLocationList(ctx, tenant, query)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Location list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listLocation
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetLocationByCode(ctx context.Context, tenant, code string) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	location, err := u.DB.GetPostgre().GetLocationByLocationCode(ctx, tenant, code)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Location")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		if err.Error() == "record not found" {
			res.Meta = helpers.GetMetaResponse(constants.RC_NOT_FOUND)
		}
		
		return res
	}
	res.Data = location
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetLocationById(ctx context.Context, tenant string, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	location, err := u.DB.GetPostgre().GetLocationByLocationId(ctx, tenant, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Location")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		if err.Error() == "record not found" {
			res.Meta = helpers.GetMetaResponse(constants.RC_NOT_FOUND)
		}
		
		return res
	}
	res.Data = location
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) DeleteLocationId(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	err := u.DB.GetPostgre().DeleteLocationById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Location")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}
