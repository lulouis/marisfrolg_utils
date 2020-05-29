package Marisfrolg_utils

import (
	"github.com/go-redis/redis"
	"time"
	"fmt"
)

func GetKeyFromRedis(key string,REDIS_CONN string) (data string, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:        REDIS_CONN,
		Password:    "", // no password set
		DB:          0,  // use default DB
		ReadTimeout: 240 * time.Second,
	})
	defer client.Close()
	data, err = client.Get(key).Result()
	if len(data) > 0 && err == nil {
		return
	} else {
		return "", fmt.Errorf(`没有找到缓存`)
	}
}

func SetKeyToRedis(keyName string,value interface{},expireTime time.Duration,REDIS_CONN string)(err error){
	client := redis.NewClient(&redis.Options{
		Addr:        REDIS_CONN,
		Password:    "", // no password set
		DB:          0,  // use default DB
		ReadTimeout: 240 * time.Second,
	})
	defer client.Close()
	err=client.Set(keyName, value, expireTime).Err()
	return
}

