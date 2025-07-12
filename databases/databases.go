package databases

import (
	"wms-server/databases/postgre"
	// "wms-server/databases/redis"
	logs "github.com/sirupsen/logrus"
)

// Database ..
type (
	database struct {
		Postgre postgre.PostgreDatabase
		// Redis redis.RedisDatabase
	}
	Database interface {
		GetPostgre() postgre.PostgreDatabase
		// GetRedis() redis.RedisDatabase
	}
)

// InitializeDatabase ..
func InitializeDatabase(
	psqlCon postgre.PostgreDatabase,
	// redisCon redis.RedisDatabase,
	l *logs.Logger,
) Database {
	return &database{
		Postgre: psqlCon,
		// Redis: redisCon,
	}
}


// GetPostgre for get connection postgre ...
func (d *database) GetPostgre() postgre.PostgreDatabase {
	return d.Postgre
}

// GetRedis for get connection redis ...
// func (d *database) GetRedis() redis.RedisDatabase {
// 	return d.Redis
// }
