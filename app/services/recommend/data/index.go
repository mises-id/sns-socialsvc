package data

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	redislib "github.com/mises-id/sns-socialsvc/lib/redis"
)

var (
	cacheSep    = ":"
	redisClient *redis.Client
	cacheKeyPre = "mises-sns"
)

type ()

func init() {
	redisClient = redislib.Client()
}

func getCacheKey(keys ...string) string {
	return cacheKeyPre + cacheSep + strings.Join(keys, cacheSep)
}

//get recommend pool union user status pool
func ListRecommendAndStarUserStatus(ctx context.Context) ([]*StatusRecommendValue, error) {
	return listRecommendedStatus(ctx)
}

func listRecommendedStatus(ctx context.Context) ([]*StatusRecommendValue, error) {
	cackeKey := getRecommendAndStatusStarUserCacheKey()
	res, err := getZAddList(ctx, cackeKey)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		keys := []string{getStatusRecommendCacheKey(), getStatusStarUserCacheKey()}
		zstore := &redis.ZStore{
			Keys: keys,
		}
		_, err = redisClient.ZUnionStore(ctx, cackeKey, zstore).Result()
		if err != nil {
			return nil, err
		}
		redisClient.Expire(ctx, cackeKey, time.Second*60*10)
		res, err = getZAddList(ctx, cackeKey)
		if err != nil {
			return nil, err
		}
	}
	return buildStatusRecommendValueSlice(res)
}

//list status recommend pool
func ListStatusRecommendPool(ctx context.Context, in *ListStatusRecommendPoolInput) ([]*StatusRecommendValue, error) {

	return listStatusRecommendPool(ctx, getStatusRecommendCacheKey())

}
func listStatusRecommendPool(ctx context.Context, cackeKey string) ([]*StatusRecommendValue, error) {

	res, err := getZAddList(ctx, cackeKey)
	if err != nil {
		return nil, err
	}
	return buildStatusRecommendValueSlice(res)
}

func getZAddList(ctx context.Context, cackeKey string) ([]string, error) {
	return redisClient.ZRangeByScore(
		ctx,
		cackeKey,
		&redis.ZRangeBy{
			Min: "-inf",
			Max: "+inf",
		},
	).Result()
}

func buildStatusRecommendValueSlice(values []string) ([]*StatusRecommendValue, error) {
	statuses := make([]*StatusRecommendValue, len(values))
	for k, v := range values {
		statuses[k] = statusRecommendValueUnMarshal(v)
	}
	return statuses, nil
}
