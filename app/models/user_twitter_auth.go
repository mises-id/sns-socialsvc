package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	TwitterUser struct {
		TwitterUserId  string    `bson:"twitter_user_id"`
		Name           string    `bson:"name"`
		UserName       string    `bson:"username"`
		FollowersCount uint64    `bson:"followers_count"`
		TweetCount     uint64    `bson:"tweet_count"`
		CreatedAt      time.Time `bson:"created_at"`
	}

	TweetInfo struct {
		TweetID   string    `bson:"tweet_id"`
		Text      string    `bson:"text"`
		CreatedAt time.Time `bson:"created_at"`
	}

	UserTwitterAuth struct {
		ID            primitive.ObjectID `bson:"_id,omitempty"`
		UID           uint64             `bson:"uid"`
		Misesid       string             `bson:"misesid,omitempty"`
		TwitterUserId string             `bson:"twitter_user_id"`
		TwitterUser   *TwitterUser       `bson:"twitter_user"`
		TweetInfo     *TweetInfo         `bson:"tweet_info"`
		CreatedAt     time.Time          `bson:"created_at"`
	}
)

func CreateUserTwitterAuthMany(ctx context.Context, data []*UserTwitterAuth) error {

	var in []interface{}
	for _, v := range data {
		in = append(in, v)
	}
	_, err := db.DB().Collection("usertwitterauths").InsertMany(ctx, in)

	return err
}

//find one user twitter auth
func FindUserTwitterAuth(ctx context.Context, params IAdminParams) (*UserTwitterAuth, error) {

	res := &UserTwitterAuth{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Last(res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

//list user twitter auth
func ListUserTwitterAuth(ctx context.Context, params IAdminParams) ([]*UserTwitterAuth, error) {

	res := make([]*UserTwitterAuth, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
func CountUserTwitterAuth(ctx context.Context, params IAdminParams) (int64, error) {

	var res int64
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Model(&UserTwitterAuth{}).Count(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}
func ListUserTwitterAuthByMisesidsOrTwitterUserIds(ctx context.Context, misesids []string, twitter_user_ids []string) ([]*UserTwitterAuth, error) {

	if len(misesids) == 0 && len(twitter_user_ids) == 0 {
		return []*UserTwitterAuth{}, nil
	}
	res := make([]*UserTwitterAuth, 0)
	chain := db.ODM(ctx).Where(bson.M{"$or": bson.A{bson.M{"misesid": bson.M{"$in": misesids}}, bson.M{"twitter_user_id": bson.M{"$in": twitter_user_ids}}}})
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
