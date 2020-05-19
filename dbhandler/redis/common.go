package redis

import (
	"github.com/vickydk/utl/log"
	"github.com/vickydk/utl/structs"
)

func ValidateData(s string) bool {
	rConn := GetRedisConnection()
	articleData, err := rConn.RDB.HGetAll(s).Result()
	if err != nil {
		log.Errorf("Redis err: %v", err)
		return false
	}
	if len(articleData) == 0 {
		return false
	}
	return true
}

func GetKeyValueByField(key string, strct interface{}) interface{} {
	rConn := GetRedisConnection()
	articleData, err := rConn.RDB.HGetAll(key).Result()
	if err != nil {
		log.Errorf("Redis err: %v", err)
		return strct
	}
	if len(articleData) == 0 {
		return strct
	}
	structs.MergeRedis(articleData, strct)
	return strct
}
