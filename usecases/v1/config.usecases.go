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

func (u *usecase) SaveConfiguration(ctx context.Context, tenant string, req cModels.RegisterConfigurationRequest) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	var configuration pModels.Configuration
	var err error
	if req.ID == 0 {
		configuration, err = u.DB.GetPostgre().GetConfigurationByConfigurationName(ctx, tenant, req.Name)
		if err != nil && err.Error() == "record not found" {

			configuration.Name = req.Name
			configuration.Value = req.Value
			configuration.Description = req.Description
			configuration.Tenant = strings.ToUpper(tenant)

			err = u.DB.GetPostgre().SaveConfiguration(ctx, configuration)
			if err != nil {
				u.Logs.WithContext(ctx).WithError(err).Error("failed save Configuration")
				res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
				return res
			}
		} else if configuration.ID != 0 {
			u.Logs.WithContext(ctx).WithError(err).Error("Configuration name already exist")
			res.Meta = helpers.GetMetaResponse(constants.RC_CONFIG_EXIST)
			return res
		}
	} else {
		configuration = pModels.Configuration{
			ID:          req.ID,
			Name:        req.Name,
			Description: req.Description,
			Value:       req.Value,
			Tenant:      strings.ToUpper(tenant),
		}

		err = u.DB.GetPostgre().SaveConfiguration(ctx, configuration)
		if err != nil {
			u.Logs.WithContext(ctx).WithError(err).Error("failed save Configuration")
			res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
			return res
		}
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetConfigurationList(ctx context.Context, tenant string, query cModels.ConfigurationQueryParams) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	listConfiguration, total, err := u.DB.GetPostgre().GetConfigurationList(ctx, tenant, query)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Configuration list")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}
	res.Data = listConfiguration
	res.Count = total
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetConfigurationByName(ctx context.Context, tenant, name string) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	Configuration, err := u.DB.GetPostgre().GetConfigurationByConfigurationName(ctx, tenant, name)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Configuration")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		if err.Error() == "record not found" {
			res.Meta = helpers.GetMetaResponse(constants.RC_NOT_FOUND)
		}
		
		return res
	}
	res.Data = Configuration
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) GetConfigurationById(ctx context.Context, tenant string, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	Configuration, err := u.DB.GetPostgre().GetConfigurationByConfigurationId(ctx, tenant, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Configuration")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		if err.Error() == "record not found" {
			res.Meta = helpers.GetMetaResponse(constants.RC_NOT_FOUND)
		}
		
		return res
	}
	res.Data = Configuration
	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}

func (u *usecase) DeleteConfigurationId(ctx context.Context, id int64) hModels.Response {
	res := hModels.Response{
		Meta: helpers.GetMetaResponse(constants.RC_GENERAL_ERROR),
	}

	err := u.DB.GetPostgre().DeleteConfigurationById(ctx, id)
	if err != nil {
		u.Logs.WithContext(ctx).WithError(err).Error("failed get Configuration")
		res.Meta = helpers.GetMetaResponse(constants.RC_GENERAL_ERROR)
		return res
	}

	res.Meta = helpers.GetMetaResponse(constants.RC_SUCCESS)
	return res
}
