package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/message"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Follow struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	FromUID        uint64             `bson:"from_uid,omitempty"`
	ToUID          uint64             `bson:"to_uid,omitempty"`
	IsFriend       bool               `bson:"is_friend,omitempty"`
	IsNew          bool               `bson:"is_new"`
	IsRead         bool               `bson:"is_read"`
	LatestPostTime *time.Time         `bson:"latest_post_time"`
	CreatedAt      time.Time          `bson:"created_at,omitempty"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty"`
	FromUser       *User              `bson:"-"`
	ToUser         *User              `bson:"-"`
}

func (a *Follow) BeforeCreate(ctx context.Context) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}

func (a *Follow) AfterCreate(ctx context.Context) error {
	_, err := CreateMessage(ctx, &CreateMessageParams{
		UID:         a.ToUID,
		FromUID:     a.FromUID,
		MessageType: enum.NewFans,
		MetaData: &message.MetaData{
			FansMeta: &message.FansMeta{
				UID: a.FromUID,
			},
		},
	})
	if err != nil {
		return err
	}
	if err = a.incrUserCounter(ctx); err != nil {
		return err
	}
	return nil
}

func (a *Follow) incrUserCounter(ctx context.Context) error {
	err := db.DB().Collection("users").FindOneAndUpdate(ctx, bson.M{"_id": a.FromUID},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   "following_count",
				Value: 1,
			}}},
		}).Err()
	if err != nil {
		return err
	}
	err = db.DB().Collection("users").FindOneAndUpdate(ctx, bson.M{"_id": a.ToUID},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   "fans_count",
				Value: 1,
			}}},
		}).Err()
	if err != nil {
		return err
	}
	return nil
}

func ReadNewFans(ctx context.Context, uid uint64) error {
	_, err := db.DB().Collection("follows").UpdateMany(ctx, bson.M{
		"to_uid": uid, "is_new": true,
	}, bson.D{{
		Key: "$set",
		Value: bson.D{{
			Key:   "is_new",
			Value: false,
		}}},
	})
	return err
}

func LatestFollowing(ctx context.Context, uid uint64) ([]*Follow, error) {
	follows := make([]*Follow, 0)
	err := db.ODM(ctx).Where(bson.M{"from_uid": uid}).
		Sort(bson.D{{Key: "is_read", Value: 1}, {Key: "latest_post_time", Value: -1}}).
		Limit(20).Find(&follows).Error
	if err != nil {
		return nil, err
	}
	return follows, preloadFollowUser(ctx, follows)
}

func ListFollow(ctx context.Context, uid uint64, relationType enum.RelationType, pageParams *pagination.QuickPagination) ([]*Follow, pagination.Pagination, error) {
	follows := make([]*Follow, 0)
	chain := db.ODM(ctx)
	if relationType == enum.Fan {
		chain = chain.Where(bson.M{"to_uid": uid})
	} else if relationType == enum.Following {
		chain = chain.Where(bson.M{"from_uid": uid})
	} else {
		chain = chain.Where(bson.M{"from_uid": uid, "is_friend": true})
	}
	paginator := pagination.NewQuickPaginator(pageParams.Limit, pageParams.NextID, chain)
	page, err := paginator.Paginate(&follows)
	if err != nil {
		return nil, nil, err
	}
	return follows, page, preloadFollowUser(ctx, follows)
}

func CreateFollow(ctx context.Context, fromUID, toUID uint64, isFriend bool) (*Follow, error) {
	follow := &Follow{
		FromUID:  fromUID,
		ToUID:    toUID,
		IsFriend: isFriend,
		IsRead:   true,
	}
	if err := follow.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	result, err := db.DB().Collection("follows").InsertOne(ctx, follow)
	if err != nil {
		return nil, err
	}

	follow.ID = result.InsertedID.(primitive.ObjectID)
	return follow, follow.AfterCreate(ctx)
}

func (f *Follow) SetFriend(ctx context.Context, isFriend bool) error {
	f.IsFriend = isFriend
	_, err := db.DB().Collection("follows").UpdateByID(ctx, f.ID, bson.M{"$set": bson.M{"is_friend": isFriend}})
	return err
}

func UpdateReadTime(ctx context.Context, uid uint64, t time.Time, targetUIDs ...uint64) error {
	return db.DB().Collection("users").FindOneAndUpdate(ctx, bson.M{
		"from_uid": uid,
		"to_uid":   bson.M{"$in": targetUIDs},
	}, bson.D{{
		Key: "$set",
		Value: bson.D{{
			Key:   "read_time",
			Value: t,
		}, {
			Key:   "updated_at",
			Value: time.Now(),
		}}},
	}).Err()
}

func GetFollow(ctx context.Context, fromUID, toUID uint64) (*Follow, error) {
	follow := &Follow{}
	result := db.DB().Collection("follows").FindOne(ctx, &bson.M{
		"from_uid": fromUID,
		"to_uid":   toUID,
	})
	err := result.Err()
	if err != nil {
		return nil, err
	}
	return follow, result.Decode(follow)
}

func EnsureDeleteFollow(ctx context.Context, fromUID, toUID uint64) error {
	_, err := GetFollow(ctx, fromUID, toUID)
	if err == nil {
		return DeleteFollow(ctx, fromUID, toUID)
	}
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

func DeleteFollow(ctx context.Context, fromUID, toUID uint64) error {
	_, err := db.DB().Collection("follows").DeleteOne(ctx, bson.M{"from_uid": fromUID, "to_uid": toUID})
	return err
}

func NewFansCount(ctx context.Context, uid uint64) (uint32, error) {
	var c int64
	return uint32(c), db.ODM(ctx).Model(&Follow{}).Where(bson.M{"to_uid": uid, "is_new": true}).Count(&c).Error
}

func UnreadUsersCount(ctx context.Context, uid uint64) (uint32, error) {
	var c int64
	return uint32(c), db.ODM(ctx).Model(&Follow{}).Where(bson.M{"from_uid": uid, "is_read": false}).Count(&c).Error
}

func ListFollowingUserIDs(ctx context.Context, uid uint64) ([]uint64, error) {
	cursor, err := db.DB().Collection("follows").Find(ctx, &bson.M{
		"from_uid": uid,
	}, &options.FindOptions{
		Projection: bson.M{"to_uid": 1},
	})
	if err != nil {
		return nil, err
	}
	follows := make([]*Follow, 0)
	if err = cursor.All(ctx, &follows); err != nil {
		return nil, err
	}
	ids := make([]uint64, len(follows))
	for i, follow := range follows {
		ids[i] = follow.ToUID
	}
	return ids, nil
}

func MarkFollowRead(ctx context.Context, fromUID, toUID uint64) error {
	t := time.Now()
	_, err := db.DB().Collection("follows").UpdateMany(ctx, bson.M{"from_uid": fromUID, "to_uid": toUID},
		bson.D{{
			Key: "$set",
			Value: bson.D{{
				Key:   "is_read",
				Value: true,
			}, {
				Key:   "updated_at",
				Value: t,
			}}},
		})
	return err
}

func preloadFollowUser(ctx context.Context, follows []*Follow) error {
	userIds := make([]uint64, 0)
	for _, follow := range follows {
		userIds = append(userIds, follow.FromUID, follow.ToUID)
	}
	users, err := FindUserByIDs(ctx, userIds...)
	if err != nil {
		return err
	}
	userMap := make(map[uint64]*User)
	for _, user := range users {
		userMap[user.UID] = user
	}
	for _, follow := range follows {
		follow.FromUser = userMap[follow.FromUID]
		follow.ToUser = userMap[follow.ToUID]
	}
	return nil
}

func GetFollowMap(ctx context.Context, fromUID uint64, toUserIDs []uint64) (map[uint64]*Follow, error) {
	follows := make([]*Follow, 0)
	err := db.ODM(ctx).Where(bson.M{
		"from_uid":   fromUID,
		"to_uid":     bson.M{"$in": toUserIDs},
		"deleted_at": nil,
	}).Find(&follows).Error
	if err != nil {
		return nil, err
	}
	followMap := make(map[uint64]*Follow)
	for _, follow := range follows {
		followMap[follow.ToUID] = follow
	}
	return followMap, nil
}
