package postgre

import (
	"context"
	"fmt"
	"strings"

	// "fmt"
	"wms-server/databases/postgre/models"
)

// FindUser ...
func (d *postgreDatabase) FindUser(ctx context.Context, email string) (models.User, error) {
	d.Logs.WithContext(ctx).WithField("email", email).Info("FindUser")
	query := d.Db.WithContext(ctx)
	user := models.User{}

	if err := query.Where("email = ? ", email).First(&user).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get user")
		return user, err
	}
	return user, nil
}

// FindUser ...
func (d *postgreDatabase) FindUserById(ctx context.Context, id int) (models.User, error) {
	d.Logs.WithContext(ctx).WithField("id", id).Info("FindUser")
	query := d.Db.WithContext(ctx)
	user := models.User{}

	if err := query.Where("id = ? ", id).First(&user).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get user")
		return user, err
	}
	return user, nil
}

// FindTenant ...
func (d *postgreDatabase) FindTenant(ctx context.Context, tenant string) (models.User, error) {
	d.Logs.WithContext(ctx).WithField("tenant", tenant).Info("FindTenant")
	query := d.Db.WithContext(ctx)
	user := models.User{}

	if err := query.Where("tenant = ?", strings.ToUpper(tenant)).First(&user).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get tenant")
		return user, err
	}
	return user, nil
}

// SaveUser ...
func (d *postgreDatabase) Save(ctx context.Context, user models.User)  error {
	d.Logs.WithContext(ctx).Info("Save User")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&user).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save user")
		return err
	}
	return nil

}


// GetList ...
func (d *postgreDatabase) GetUserList(ctx context.Context,tenant, email string, page , limit int) ( []models.User, int64, error) {
	query := d.Db.WithContext(ctx)
	
	var res []models.User
	var total int64
	fmt.Println("ini tenant ", tenant)
	query = query.Where("tenant = ?", tenant)
	if email != "" {
		query = query.Where("email = ?", strings.ToLower(email))
	}

	err := query.Order("id asc").Limit(limit).Offset((page - 1) * limit).Find(&res).Limit(-1).Offset(-1).Count(&total).Error

	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error get Customers list")
		return res, 0, err
	}

	return res, total, nil
}

// Ping ...
func (d *postgreDatabase) Ping(ctx context.Context) error {
	sql, err := d.Db.WithContext(ctx).DB()
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error connecting DB")
		return err
	}
	err = sql.Ping()
	if err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error ping database")
		return err
	}
	return nil
}

