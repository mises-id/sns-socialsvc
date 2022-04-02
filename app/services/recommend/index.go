package recommend

import (
	"strings"

	"github.com/go-redis/redis/v8"
	redislib "github.com/mises-id/sns-socialsvc/lib/redis"
)

var (
	cacheSep    = ":"
	redisClient *redis.Client
	cacheKeyPre = "mises-sns"
)

func init() {
	redisClient = redislib.Client()
}

func getCacheKey(keys ...string) string {
	return cacheKeyPre + cacheSep + strings.Join(keys, cacheSep)
}
