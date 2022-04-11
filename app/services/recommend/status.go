package recommend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mises-id/sns-socialsvc/admin"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	statusKey          = "status"
	statusRecommendKey = "recommend"
	statusGroupUserKey = "group_user"
	statusLastID       primitive.ObjectID
)

type (
	StatusGroupUserPool struct {
		ctx      context.Context
		cackeKey string
		uid      uint64
	}
	StatusRecommendPool struct {
		ctx      context.Context
		cackeKey string
	}
	StatusRecommendValue struct {
		ID  primitive.ObjectID `json:"id"`
		UID uint64             `json:"uid"`
	}
)

func getStatusCacheKey(key string) string {
	return getCacheKey([]string{statusKey, key}...)
}
func getStatusRecommendCacheKey() string {
	return getStatusCacheKey(statusRecommendKey)
}
func getStatusGroupUserCacheKey(uid uint64) string {
	return getStatusCacheKey(statusGroupUserKey + cacheSep + strconv.Itoa(int(uid)))
}

//status group user pool
func NewStatusGroupUserPool(ctx context.Context, uid uint64) *StatusGroupUserPool {
	return &StatusGroupUserPool{
		ctx:      ctx,
		uid:      uid,
		cackeKey: getStatusGroupUserCacheKey(uid),
	}
}

func StatusGroupUserPoolInit(ctx context.Context) error {
	c, err := models.CountStatus(ctx, &admin.AdminStatusParams{})
	if err != nil {
		fmt.Println("count status error: ", err.Error())
		return err
	}
	if c == 0 {
		return nil
	}
	var listNum int64
	listNum = 50
	times := int(math.Ceil(float64(c) / float64(listNum)))
	for i := 0; i < times; i++ {
		err := initStatusGroupUserPool(ctx, listNum)
		if err != nil {
			fmt.Println("do error: ", err.Error())
			return err
		}
	}

	return nil
}

func initStatusGroupUserPool(ctx context.Context, num int64) error {

	statuses, err := models.AdminListStatus(ctx, &admin.AdminStatusParams{ListNum: num, LastID: statusLastID})
	if err != nil {
		if err != nil {
			fmt.Println("find error: ", err.Error())
		}
		return err
	}
	for _, status := range statuses {
		err := addStatusGroupUserPool(ctx, getStatusGroupUserCacheKey(status.UID), status)
		if err != nil {
			fmt.Println("add error: ", err.Error())
		}
	}
	statusLastID = statuses[len(statuses)-1].ID
	return nil

}

func AddStatusGroupUserPool(ctx context.Context, status *models.Status) error {
	uid := status.UID
	return addStatusGroupUserPool(ctx, getStatusGroupUserCacheKey(uid), status)
}

func addStatusGroupUserPool(ctx context.Context, cackeKey string, status *models.Status) error {
	meb := &redis.Z{
		Score:  float64(status.CreatedAt.UnixMilli()),
		Member: status.ID.Hex(),
	}
	return redisClient.ZAdd(ctx, cackeKey, meb).Err()
}

func RemoveStatusGroupUserPool(ctx context.Context, uid uint64, statusid primitive.ObjectID) error {

	return removeStatusGroupUserPool(ctx, getStatusGroupUserCacheKey(uid), statusid)

}
func removeStatusGroupUserPool(ctx context.Context, cackeKey string, statusid primitive.ObjectID) error {

	return redisClient.ZRem(ctx, cackeKey, statusid).Err()
}

// recommend pool
func NewStatusRecommendPool(ctx context.Context) *StatusRecommendPool {
	return &StatusRecommendPool{
		ctx:      ctx,
		cackeKey: getStatusRecommendCacheKey(),
	}
}

func NewStatusRecommendInit(ctx context.Context) error {

	c, err := models.CountStatus(ctx, &admin.AdminStatusParams{
		Tag:      enum.TagRecommendStatus,
		OnlyShow: true,
	})
	if err != nil {
		fmt.Println("count status error: ", err.Error())
		return err
	}
	if c == 0 {
		return nil
	}
	var listNum int64
	listNum = 50
	last_id := primitive.NilObjectID
	ctx = context.WithValue(ctx, primitive.ObjectID{}, last_id)
	times := int(math.Ceil(float64(c) / float64(listNum)))
	for i := 0; i < times; i++ {
		err := initStatusGroupUserPool(ctx, listNum)
		if err != nil {
			fmt.Println("do error: ", err.Error())
			return err
		}
	}
	return nil
}

//init recommend pool
func initStatusRecommendPool(ctx context.Context, num int64) error {
	last_id, ok := ctx.Value(primitive.ObjectID{}).(primitive.ObjectID)
	if !ok {
		return errors.New("no last_id")
	}
	statuses, err := models.AdminListStatus(ctx, &admin.AdminStatusParams{
		ListNum:  num,
		LastID:   last_id,
		Tag:      enum.TagRecommendStatus,
		OnlyShow: true,
	})
	if err != nil {
		if err != nil {
			fmt.Println("find error: ", err.Error())
		}
		return err
	}
	for _, status := range statuses {

		err := addStatusRecommendPool(ctx, getStatusRecommendCacheKey(), status)
		if err != nil {
			fmt.Println("add error: ", err.Error())
		}
	}
	last_id = statuses[len(statuses)-1].ID
	ctx = context.WithValue(ctx, primitive.ObjectID{}, last_id)
	return nil
}

func AddStatusRecommendPool(ctx context.Context, status *models.Status) error {

	return addStatusRecommendPool(ctx, getStatusRecommendCacheKey(), status)

}

func addStatusRecommendPool(ctx context.Context, cackeKey string, status *models.Status) error {
	meb := &redis.Z{
		Score:  float64(time.Now().UnixMilli()),
		Member: statusToRecommendValue(status),
	}
	return redisClient.ZAdd(ctx, cackeKey, meb).Err()
}

func RemoveStatusRecommendPool(ctx context.Context, status *models.Status) error {

	return removeStatusRecommendPool(ctx, getStatusRecommendCacheKey(), status)

}

func removeStatusRecommendPool(ctx context.Context, cackeKey string, status *models.Status) error {
	return redisClient.ZRem(ctx, cackeKey, statusToRecommendValue(status)).Err()
}

func statusToRecommendValue(status *models.Status) []byte {
	str, _ := json.Marshal(StatusRecommendValue{ID: status.ID, UID: status.UID})
	return str
}
