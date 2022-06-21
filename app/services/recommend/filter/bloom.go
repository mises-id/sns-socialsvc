package filter

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	statusBloomCacheKey = "status-bloom"
)

type ()

func getUserStatusBloomCacheKey(uid uint64) string {
	return getCacheKey(statusBloomCacheKey + cacheSep + strconv.Itoa(int(uid)))
}

//user status
func StatusBfInsert(ctx context.Context, uid uint64, statusids ...primitive.ObjectID) error {

	return statusBfInsert(ctx, uid, statusids...)
}

func statusBfInsert(ctx context.Context, uid uint64, statusids ...primitive.ObjectID) error {
	var id interface{}
	cacheKey := getUserStatusBloomCacheKey(uid)
	values := make([]interface{}, len(statusids))
	for k, statusid := range statusids {
		id = statusid.Hex()
		values[k] = id
	}
	op := &redis.BfInsertArgs{
		Key:    cacheKey,
		Values: values,
	}
	return redisClient.BfInsert(ctx, op).Err()

}

//user statusid bf exists

func StatusBfMexists(ctx context.Context, uid uint64, statusids ...primitive.ObjectID) ([]bool, error) {

	return statusBfMexists(ctx, uid, statusids...)
}

func StatusBfExists(ctx context.Context, uid uint64, statusid primitive.ObjectID) (bool, error) {

	var id interface{}
	id = statusid.Hex()
	return redisClient.BfExists(ctx, getUserStatusBloomCacheKey(uid), id).Bool()
}

func statusBfMexists(ctx context.Context, uid uint64, statusids ...primitive.ObjectID) ([]bool, error) {
	var v interface{}
	cacheKey := getUserStatusBloomCacheKey(uid)
	values := make([]interface{}, len(statusids))
	for k, id := range statusids {
		v = id.Hex()
		values[k] = v
	}
	return redisClient.BfMexists(ctx, cacheKey, values).BoolSlice()

}
