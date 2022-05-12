package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
)

func init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       6,  // use default DB
		//OnConnect: con,
	})
	redisClient = rdb
}

func con(ctx context.Context, cn *redis.Conn) error {
	return errors.New("on connect")
}

func SetupRedis(ctx context.Context) {

}

func Client() *redis.Client {
	return redisClient
}
