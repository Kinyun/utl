package redis

import (
	"github.com/go-redis/redis"
	"github.com/vickydk/utl/log"
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

func GetRedisConnection(RedisAddress, RedisPassword string) *connections {
	if connection == nil {
		lockconnection.Lock()
		defer lockconnection.Unlock()
		connection = newConnection(RedisAddress, RedisPassword)
	}

	return connection
}

func newConnection(RedisAddress, RedisPassword string) *connections {
	rdb := redis.NewClient(&redis.Options{
		Addr:         RedisAddress,
		Password:     RedisPassword,
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
