package redisutil

import (
	"os"
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

var once = sync.Once{}
var rdb *redis.Client

func GetRedisClient() *redis.Client {
	once.Do(func() {
		db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			db = 0
		}
		rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASS"), // no password set
			DB:       db,                      // use default DB
		})
	})
	return rdb
}
