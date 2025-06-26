package redis

import (
	"context"
	"strconv"
	"time"

	"wms-server/helpers"
	uModels "wms-server/usecases/v1/models"

	redis "github.com/redis/go-redis/v9"
	logs "github.com/sirupsen/logrus"
)

type (
	redisDatabase struct {
		ClientRed *redis.Client
		Logs      *logs.Logger
	}
	// RedisDatabase ...
	RedisDatabase interface {
		GetQueue(ctx context.Context, d time.Duration, key string) ([]string, error)
		PushJobRedis(ctx context.Context, key string, value interface{}) error
		SaveRedisExp(ctx context.Context, exp time.Duration, key string, val interface{}) error
		GetRedisKey(ctx context.Context, Key string) (string, error)
		HeatchCheck(ctx context.Context) uModels.DataHealthCheck
		Increment(ctx context.Context, key string) error
	}
)

// InitializeRedis ..
func InitializeRedis(conn *redis.Client, log *logs.Logger) RedisDatabase {
	return &redisDatabase{
		ClientRed: conn,
		Logs:      log,
	}
}

// ConnectRedis ...
func ConnectRedis() *redis.Client {
	c := context.Background()
	address := helpers.GetEnv("REDIS_ADDRESS")
	port := helpers.GetEnv("REDIS_PORT")
	dbtype, _ := strconv.Atoi(helpers.GetEnv("REDIS_DB_TYPE"))
	password := helpers.GetEnv("REDIS_PASSWORD")

	ClientRed := redis.NewClient(&redis.Options{
		Addr:     address + ":" + port,
		Password: password, // no password set
		DB:       dbtype,   // use default DB
	})

	_, err := ClientRed.Ping(c).Result()
	if err != nil {
		logs.Error("Redis Err ", err)
		panic("Error open redis connection")
	}

	logs.Info("Redis connected successfully")
	return ClientRed

}
