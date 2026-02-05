package postgre

import (
	"context"
	cModels "wms-server/controllers/v1/models"
	"wms-server/databases/postgre/models"
)


func (d *postgreDatabase) SaveAccountPayable(ctx context.Context, data models.AccountPayable) error {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveAccountPayable")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save AccountPayable")
		return err
	}
	return nil

}


func (d *postgreDatabase) GetAccountPayableList(ctx context.Context, tenant string, req cModels.ReqListFinance) ([]models.AccountPayable, int64, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.AccountPayable
	var total int64

	query = query.Where("tenant = ?", tenant)

	query = query.Where("tenant = ?", tenant)
	if req.PaymentMethod != "" {
		query = query.Where("payment_method = ?", req.PaymentMethod)
	}

	if req.ReferenceNumber != "" {
		query = query.Where("reference_number = ?", req.ReferenceNumber)
	}

	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("payment_date between ? and ?", req.StartDate + " 00:00:00" , req.EndDate + " 23:59:59")
	}
	
	err := query.Order("id desc").Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get AccountPayable list")
		return res, 0, err
	}

	return res, total, nil

}

func (d *postgreDatabase) SaveAccountReceivable(ctx context.Context, data models.AccountReceiveble) error {
	d.Logs.WithContext(ctx).WithField("data", data).Info("SaveAccountReceivable")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&data).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save AccountReceivable")
		return err
	}
	return nil

}

func (d *postgreDatabase) GetAccountReceivableList(ctx context.Context, tenant string, req cModels.ReqListFinance) ([]models.AccountReceiveble, int64, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.AccountReceiveble
	var total int64

	query = query.Where("tenant = ?", tenant)
	if req.PaymentMethod != "" {
		query = query.Where("payment_method = ?", req.PaymentMethod)
	}

	if req.ReferenceNumber != "" {
		query = query.Where("reference_number = ?", req.ReferenceNumber)
	}

	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("payment_date between ? and ?", req.StartDate + " 00:00:00" , req.EndDate + " 23:59:59")
	}

	err := query.Order("id desc").Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get AccountReceivable List")
		return res, 0, err
	}

	return res, total, nil

}

func (d *postgreDatabase) GetAccountReceivableTotal(ctx context.Context, tenant string, req cModels.ReqListFinance) (float64, error) {
	query := d.Db.WithContext(ctx)
	var total float64
	query = query.Where("tenant = ?", tenant)
	if req.ReferenceNumber != "" {
		query = query.Where("reference_number = ?", req.ReferenceNumber)
	}
	
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("payment_date between ? and ?", req.StartDate + " 00:00:00" , req.EndDate + " 23:59:59")
	}

	err := query.Table("account_receivable ap").Select("sum(ap.amount ) as total").Scan(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get AccountReceivable Total")
		return 0, err
	}

	return total, err
}

func (d *postgreDatabase) GetAccountPayableTotal(ctx context.Context, tenant string, req cModels.ReqListFinance) (float64, error) {
		query := d.Db.WithContext(ctx)
	var total float64
	
	query = query.Where("tenant = ?", tenant)
	if req.ReferenceNumber != "" {
		query = query.Where("reference_number = ?", req.ReferenceNumber)
	}
	
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("payment_date between ? and ?", req.StartDate , req.EndDate )
	}

	err := query.Table("account_payable ap").Select("sum(ap.amount ) as total").Scan(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get AccountPayableTotal")
		return 0, err
	}

	return total, err
}