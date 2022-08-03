package data

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	userKey          = "user"
	userFollowingKey = "following"
	userFollowLastID primitive.ObjectID
)

type ()

func getUserCacheKey(key string) string {
	return getCacheKey([]string{userKey, key}...)
}

func getUserFollowingCacheKey(uid uint64) string {
	return getUserCacheKey(userFollowingKey + cacheSep + strconv.Itoa(int(uid)))
}

//===================
//user following pool
//===================
func InitUserFollowingPool(ctx context.Context) error {

	c, err := models.CountFollow(ctx, &search.FollowSearch{})
	if err != nil {
		fmt.Println("count status error: ", err.Error())
		return err
	}
	if c == 0 {
		return nil
	}
	fmt.Println("follow count: ", c)
	var listNum int64
	listNum = 50
	times := int(math.Ceil(float64(c) / float64(listNum)))
	for i := 0; i < times; i++ {
		err := initUserFollowingPool(ctx, listNum)
		if err != nil {
			fmt.Println("do error: ", err.Error())
			return err
		}
	}

	return nil
}
func initUserFollowingPool(ctx context.Context, num int64) error {

	follows, err := models.NewListFollow(ctx, &search.FollowSearch{ListNum: num, LastID: userFollowLastID})
	if err != nil {
		if err != nil {
			fmt.Println("find error: ", err.Error())
		}
		return err
	}
	for _, follow := range follows {
		err := addUserFollowingPool(ctx, follow.FromUID, follow.ToUID)
		if err != nil {
			fmt.Println("add error: ", err.Error())
		}
	}
	userFollowLastID = follows[len(follows)-1].ID
	return nil

}

//add following
func AddUserFollowingPool(ctx context.Context, fromUID, toUID uint64) error {

	return addUserFollowingPool(ctx, fromUID, toUID)
}

func addUserFollowingPool(ctx context.Context, fromUID, toUID uint64) error {
	cackeKey := getUserFollowingCacheKey(fromUID)
	meb := &redis.Z{
		Score:  float64(time.Now().UnixMilli()),
		Member: toUID,
	}
	return redisClient.ZAdd(ctx, cackeKey, meb).Err()
}

//remove following
func RemoveUserFollowingPool(ctx context.Context, fromUID, toUID uint64) error {

	return removeUserFollowingPool(ctx, fromUID, toUID)

}
func removeUserFollowingPool(ctx context.Context, fromUID, toUID uint64) error {
	cackeKey := getUserFollowingCacheKey(fromUID)
	return redisClient.ZRem(ctx, cackeKey, toUID).Err()
}

func getFollow2User(ctx context.Context, uid uint64) ([]uint64, error) {
	if uid <= 0 {
		return []uint64{}, nil
	}
	//user_followings,err :=
	return nil, nil
}
