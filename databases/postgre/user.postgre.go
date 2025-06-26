package postgre

import (
	"context"
	"strings"
	// "fmt"
	"wms-server/databases/postgre/models"
)

// FindUser ...
func (d *postgreDatabase) FindUser(ctx context.Context, email string) (models.User, error) {
	d.Logs.WithContext(ctx).WithField("email", email).Info("FindUser")
	query := d.Db.WithContext(ctx)
	user := models.User{}

	if err := query.Where("email = ?", email).First(&user).Error; err != nil {
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

// FindTenant ...
func (d *postgreDatabase) Save(ctx context.Context, user models.User)  error {
	d.Logs.WithContext(ctx).Info("Save User")
	query := d.Db.WithContext(ctx)

	if err := query.Save(&user).Error; err != nil {
		d.Logs.WithContext(ctx).WithError(err).Error("Error save user")
		return err
	}
	return nil

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

