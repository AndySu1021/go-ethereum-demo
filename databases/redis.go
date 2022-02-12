package databases

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     getRedisAddr(),
		Password: viper.GetString("databases.redis.password"),
		DB:       0,
	})
}

func getRedisAddr() string {
	return viper.GetString("databases.redis.host") + ":" + viper.GetString("databases.redis.port")
}
