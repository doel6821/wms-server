package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"wms-server/helpers"
	"time"

	uModels "wms-server/usecases/v1/models"

	logs "github.com/sirupsen/logrus"
)

// GetQueue redis ...
func (r *redisDatabase) GetQueue(ctx context.Context, d time.Duration, key string) ([]string, error) {
	start := time.Now()
	defer func() {
		r.Logs.WithContext(ctx).WithFields(logs.Fields{
			"OPERATION": "QUEUE",
			"KEY":       key,
			"DURATION":  time.Since(start).Milliseconds(),
		}).Info()
	}()
	return r.ClientRed.BLPop(ctx, 0*time.Second, key).Result()
}

// PushJobRedis redis ..
func (r *redisDatabase) PushJobRedis(ctx context.Context, key string, value interface{}) error {
	start := time.Now()
	defer func() {
		r.Logs.WithContext(ctx).WithFields(logs.Fields{
			"OPERATION": "RPUSH",
			"KEY":       key,
			"DURATION":  time.Since(start).Milliseconds(),
		}).Info()
	}()
	var data []byte
	var err error

	data, err = json.Marshal(value)
	if err != nil {
		r.Logs.WithContext(ctx).WithError(err).Error("Error marshalling")
		return err
	}
	for i := 0; i < 3; i++ {
		if err := r.ClientRed.RPush(ctx, key, string(data)).Err(); err == nil {
			break
		}
		r.Logs.WithContext(ctx).WithError(err).Error("Error Push")
	}

	return err
}

// SaveRedisExp ..
func (r *redisDatabase) SaveRedisExp(ctx context.Context, exp time.Duration, key string, val interface{}) error {
	start := time.Now()
	defer func() {
		r.Logs.WithContext(ctx).WithFields(logs.Fields{
			"OPERATION": "SET",
			"KEY":       key,
			"DURATION":  time.Since(start).Milliseconds(),
		}).Info()
	}()
	var err error
	for i := 0; i < 3; i++ {
		err = r.ClientRed.Set(ctx, key, val, exp).Err()
		if err == nil {
			break
		}
		r.Logs.WithContext(ctx).WithError(err).Error("Error Set")
	}
	return err
}

// GetRedisKey ..
func (r *redisDatabase) GetRedisKey(ctx context.Context, key string) (string, error) {
	start := time.Now()
	defer func() {
		r.Logs.WithContext(ctx).WithFields(logs.Fields{
			"OPERATION": "GET",
			"KEY":       key,
			"DURATION":  time.Since(start).Milliseconds(),
		}).Info()
	}()
	value, err := r.ClientRed.Get(ctx, key).Result()
	if err != nil {
		r.Logs.WithContext(ctx).WithError(err).Error("Error get")
		return value, err
	}
	return value, err
}

// HeatchCheck ...
func (r *redisDatabase) HeatchCheck(ctx context.Context) uModels.DataHealthCheck {
	res := uModels.DataHealthCheck{
		ServiceName: "Redis Single",
		Host:        helpers.GetEnv("REDIS_ADDRESS"),
		StatusCode:  200,
	}

	_, err := r.ClientRed.Ping(ctx).Result()
	if err != nil {
		r.Logs.WithContext(ctx).WithError(err).Error("Error ping")
		res.StatusCode = 400
		return res
	}
	info, err := r.ClientRed.Do(ctx, "INFO", "server").Result()
	res.AdditionalData = fmt.Sprintf("%v, err : %v", info, err)
	return res
}

// Increment ...
func (r *redisDatabase) Increment(ctx context.Context, key string) error {
	start := time.Now()
	defer func() {
		r.Logs.WithContext(ctx).WithFields(logs.Fields{
			"OPERATION": "INCR",
			"KEY":       key,
			"DURATION":  time.Since(start).Milliseconds(),
		}).Info()
	}()
	return r.ClientRed.Incr(ctx, key).Err()
}

// Expire ...
func (r *redisDatabase) Expire(ctx context.Context, key string, d time.Duration) error {
	start := time.Now()
	defer func() {
		r.Logs.WithContext(ctx).WithFields(logs.Fields{
			"OPERATION": "EXPIRE",
			"KEY":       key,
			"DURATION":  time.Since(start).Milliseconds(),
		}).Info()
	}()
	return r.ClientRed.Expire(ctx, key, d).Err()
}
