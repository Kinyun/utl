package redis

import (
	"github.com/go-redis/redis"
	"omsApi/pkg/utl/config"
	"omsApi/pkg/utl/log"
	"sync"
	"time"
)

type connections struct {
	RDB *redis.Client
}

var (
	connection     *connections
	lockconnection = &sync.Mutex{}
)

func GetRedisConnection() *connections {
	if connection == nil {
		lockconnection.Lock()
		defer lockconnection.Unlock()
		connection = newConnection()
	}

	return connection
}

func newConnection() *connections {
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.Env.RedisAddress,
		Password:     config.Env.RedisPassword,
		PoolTimeout:  20 * time.Second,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Error(err)
	}

	return &connections{RDB: rdb}
}
