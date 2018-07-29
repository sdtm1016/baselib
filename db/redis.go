package db

import (
	"github.com/go-redis/redis"
	"baselib/logger"
	"baselib/config"
)

var client *redis.Client

func InitRedis() error {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     config.GetConf().RedisDb.Addr,
			Password: config.GetConf().RedisDb.Password,
			DB:       config.GetConf().RedisDb.Database,
		})
	}

	_, err := client.Ping().Result()
	if err != nil {
		logger.Error("fail to connect redis, error:", err)
		return err
	}

	logger.Info("succeed connect to redis")
	return nil
}

func GetRedisClient() *redis.Client {
	if client == nil {
		InitRedis()
	}
	return client
}
