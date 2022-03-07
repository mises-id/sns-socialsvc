package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Following2PoolCursor struct {
		Max int64 `bson:"max,omitempty"`
		Min int64 `bson:"min,omitempty"`
	}
	RecommendStatusPoolCursor struct {
		Max int64 `bson:"max,omitempty"`
		Min int64 `bson:"min,omitempty"`
	}
	CommonPoolCursor struct {
		Max int64 `bson:"max,omitempty"`
		Min int64 `bson:"min,omitempty"`
	}

	UserExt struct {
		ID                        primitive.ObjectID         `bson:"_id,omitempty"`
		UID                       uint64                     `bson:"uid"`
		AirdropCoin               uint64                     `bson:"airdrop_coin"`
		IsLogined                 bool                       `bson:"is_logined"`
		TwitterAirdrop            bool                       `bson:"twitter_airdrop"`
		LastViewTime              time.Time                  `bson:"last_view_time"`
		RecommendStatusPoolCursor *RecommendStatusPoolCursor `bson:"recommend_status_pool_cursor"`
		Following2PoolCursor      *Following2PoolCursor      `bson:"following2_cursor"`
		CommonPoolCursor          *CommonPoolCursor          `bson:"common_cursor"`
		CreatedAt                 time.Time                  `bson:"created_at,omitempty"`
		UpdatedAt                 time.Time                  `bson:"updated_at,omitempty"`
	}
)

//find or create user ext
func FindOrCreateUserExt(ctx context.Context, uid uint64) (*UserExt, error) {

	//find
	user_ext := &UserExt{}
	err := db.ODM(ctx).Where(bson.M{
		"uid": uid,
	}).Last(user_ext).Error

	if err == mongo.ErrNoDocuments {
		return CreateUserExt(ctx, uid)
	}
	if err != nil {
		return nil, err
	}
	return user_ext, nil

}

func UserMergeUserExt(ctx context.Context, user *User) *User {
	user_ext, err := FindOrCreateUserExt(ctx, user.UID)
	if err != nil {
		return user
	}
	user.IsAirdropped = user_ext.TwitterAirdrop
	user.IsLogined = user_ext.IsLogined
	if !user_ext.IsLogined {
		user_ext.updateIsLogin(ctx)
	}
	return user
}

func (m *UserExt) BeforeSave(ctx context.Context) error {

	//create
	if m.ID == primitive.NilObjectID {
		m.CreatedAt = time.Now()
	}
	m.UpdatedAt = time.Now()
	return nil
}

func (m *UserExt) UpdateAirdrop(ctx context.Context) error {
	update := bson.M{}
	update["twitter_airdrop"] = true
	update["airdrop_coin"] = m.AirdropCoin
	_, err := db.DB().Collection("userexts").UpdateOne(ctx, &bson.M{
		"uid": m.UID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}
func (m *UserExt) updateIsLogin(ctx context.Context) error {
	update := bson.M{}
	update["is_logined"] = true
	_, err := db.DB().Collection("userexts").UpdateOne(ctx, &bson.M{
		"uid": m.UID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err
}

//update user ext
func (m *UserExt) Update(ctx context.Context) error {

	/* err := m.BeforeSave(ctx)
	if err != nil {
		return err
	} */
	nt := time.Now()
	update := bson.M{}
	update["last_view_time"] = nt
	update["updated_at"] = nt
	if m.RecommendStatusPoolCursor != nil {
		update["recommend_status_pool_cursor"] = m.RecommendStatusPoolCursor
	}
	if m.Following2PoolCursor != nil {
		update["following2_cursor"] = m.Following2PoolCursor
	}
	if m.CommonPoolCursor != nil {
		update["common_cursor"] = m.CommonPoolCursor
	}
	_, err := db.DB().Collection("userexts").UpdateOne(ctx, &bson.M{
		"uid": m.UID,
	}, bson.D{{
		Key:   "$set",
		Value: update}})
	return err

}

//create user ext
func CreateUserExt(ctx context.Context, uid uint64) (*UserExt, error) {
	user_ext := &UserExt{
		UID: uid,
	}
	if err := user_ext.BeforeSave(ctx); err != nil {
		return nil, err
	}
	if err := db.ODM(ctx).Create(user_ext).Error; err != nil {
		return nil, err
	}

	return user_ext, nil
}
