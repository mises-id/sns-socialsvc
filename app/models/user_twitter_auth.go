package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	TweetInfo struct {
		TweetID   string    `bson:"tweet_id"`
		Text      string    `bson:"text"`
		CreatedAt time.Time `bson:"created_at"`
	}
	TwitterUser struct {
		TwitterUserId  string    `bson:"twitter_user_id"`
		Name           string    `bson:"name"`
		UserName       string    `bson:"username"`
		FollowersCount uint64    `bson:"followers_count"`
		FollowingCount uint64    `bson:"following_count"`
		TweetCount     uint64    `bson:"tweet_count"`
		CreatedAt      time.Time `bson:"created_at"`
	}
	CheckResult struct {
		CheckNum          int `bson:"check_num"`
		ZeroTweetNum      int `bson:"zero_tweet_num"`
		ZeroFollowerNum   int `bson:"zero_follower_num"`
		TotalFollowerNum  int `bson:"total_follower_num"`
		RecentRegisterNum int `bson:"recent_register_num"`
	}
	UserTwitterAuth struct {
		ID                   primitive.ObjectID `bson:"_id,omitempty"`
		UID                  uint64             `bson:"uid"`
		Misesid              string             `bson:"misesid,omitempty"`
		TwitterUserId        string             `bson:"twitter_user_id"`
		OauthToken           string             `bson:"oauth_token"`
		OauthTokenSecret     string             `bson:"oauth_token_secret"`
		TwitterUser          *TwitterUser       `bson:"twitter_user"`
		CheckResult          *CheckResult       `bson:"check_result"`
		TweetInfo            *TweetInfo         `bson:"tweet_info"`
		UpdatedAt            time.Time          `bson:"updated_at,omitempty"`
		CreatedAt            time.Time          `bson:"created_at"`
		Amount               int64              `bson:"amount"`
		IsValid              bool               `bson:"-"`
		IsAirdrop            bool               `bson:"is_airdrop"`
		SendTweeState        int                `bson:"send_tweet_state"`        // 1 pending 2 success 3 failed 4 Unauthorized
		LikeTweeState        int                `bson:"like_tweet_state"`        // 1 pending 2 success 3 failed 4 Unauthorized
		FindTwitterUserState int                `bson:"find_twitter_user_state"` // 1 pending 2 success 3 failed 4 Unauthorized
		FollowState          int                `bson:"follow_state"`            // 1 pending 2 success 3 failed 4 Unauthorized
		ValidState           int                `bson:"valid_state"`             // 2 valid 3 invalid 4 need check
		IsFollowed           bool               `bson:"is_followed"`
		Ipaddr               string             `bson:"ipaddr"`
	}
)

func CreateUserTwitterAuth(ctx context.Context, data *UserTwitterAuth) error {
	created := bson.M{}
	created["created_at"] = time.Now()
	created["updated_at"] = time.Now()
	created["uid"] = data.UID
	created["misesid"] = data.Misesid
	created["twitter_user_id"] = data.TwitterUserId
	created["twitter_user"] = data.TwitterUser
	created["tweet_info"] = data.TweetInfo
	created["oauth_token"] = data.OauthToken
	created["oauth_token_secret"] = data.OauthTokenSecret
	created["find_twitter_user_state"] = data.FindTwitterUserState
	opt := &options.FindOneAndUpdateOptions{}
	opt.SetUpsert(true)
	opt.SetReturnDocument(1)
	result := db.DB().Collection("usertwitterauths").FindOneAndUpdate(ctx,
		bson.M{"uid": data.UID, "twitter_user_id": data.TwitterUserId},
		bson.D{{Key: "$set", Value: created}}, opt)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func UpdateUserTwitterAuth(ctx context.Context, data *UserTwitterAuth) error {

	update := bson.M{}
	update["updated_at"] = time.Now()
	update["twitter_user_id"] = data.TwitterUserId
	update["find_twitter_user_state"] = data.FindTwitterUserState
	update["twitter_user"] = data.TwitterUser
	update["oauth_token"] = data.OauthToken
	update["oauth_token_secret"] = data.OauthTokenSecret

	_, err := db.DB().Collection("usertwitterauths").UpdateByID(ctx, data.ID, bson.D{{Key: "$set", Value: update}})
	return err
}
func UpdateUserTwitterAuthFollow(ctx context.Context, data *UserTwitterAuth) error {

	update := bson.M{}
	update["follow_state"] = data.FollowState
	_, err := db.DB().Collection("usertwitterauths").UpdateByID(ctx, data.ID, bson.D{{Key: "$set", Value: update}})
	return err
}
func UpdateUserTwitterAuthSendTweet(ctx context.Context, data *UserTwitterAuth) error {

	update := bson.M{}
	update["send_tweet_state"] = data.SendTweeState
	update["like_tweet_state"] = data.LikeTweeState
	_, err := db.DB().Collection("usertwitterauths").UpdateByID(ctx, data.ID, bson.D{{Key: "$set", Value: update}})
	return err
}
func UpdateUserTwitterAuthFindState(ctx context.Context, data *UserTwitterAuth) error {

	update := bson.M{}
	update["find_twitter_user_state"] = data.FindTwitterUserState
	_, err := db.DB().Collection("usertwitterauths").UpdateByID(ctx, data.ID, bson.D{{Key: "$set", Value: update}})
	return err
}

func DeleteUserTwitterAuthByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := db.DB().Collection("usertwitterauths").DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func UpdateUserTwitterAuthTwitterUser(ctx context.Context, data *UserTwitterAuth) error {

	update := bson.M{}
	update["twitter_user"] = data.TwitterUser
	update["amount"] = data.Amount
	update["follow_state"] = data.FollowState
	update["valid_state"] = data.ValidState
	update["is_airdrop"] = data.IsAirdrop
	update["find_twitter_user_state"] = data.FindTwitterUserState
	update["send_tweet_state"] = data.SendTweeState
	if data.CheckResult != nil {
		update["check_result"] = data.CheckResult
	}
	_, err := db.DB().Collection("usertwitterauths").UpdateByID(ctx, data.ID, bson.D{{Key: "$set", Value: update}})
	return err
}

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
func FindUserTwitterAuthByUid(ctx context.Context, uid uint64) (*UserTwitterAuth, error) {
	res := &UserTwitterAuth{}
	err := db.ODM(ctx).Where(bson.M{"uid": uid}).Last(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
func FindUserTwitterAuthByTwitterUserId(ctx context.Context, twitter_user_id string) (*UserTwitterAuth, error) {
	res := &UserTwitterAuth{}
	err := db.ODM(ctx).Where(bson.M{"twitter_user_id": twitter_user_id}).Last(&res).Error
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
