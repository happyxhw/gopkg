package goredis

import (
	"github.com/go-redis/redis/v7"
)

// RedisConn redis client
//var pool *redis.Pool

type Config struct {
	Host         string
	Password     string
	Db           int
	PoolSize     int `mapstructure:"pool_size"`
	MinIdleConns int `mapstructure:"min_idle_conns"`
}

// NewRedis Initialize the Redis instance
func NewRedis(redisConf *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisConf.Host,
		DB:           redisConf.Db,
		Password:     redisConf.Password,
		PoolSize:     redisConf.PoolSize,
		MinIdleConns: redisConf.MinIdleConns,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return client, nil
}
