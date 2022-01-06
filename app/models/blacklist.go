package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Blacklist struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UID       uint64             `bson:"uid,omitempty"`
	TargetUID uint64             `bson:"target_uid,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

func ListBlacklist(ctx context.Context, uid uint64) ([]*Blacklist, error) {
	blacklist := make([]*Blacklist, 0)
	return blacklist, db.ODM(ctx).Where(bson.M{"uid": uid}).Find(&blacklist).Error
}

func FindBlacklist(ctx context.Context, uid, targetUID uint64) (*Blacklist, error) {
	blacklist := &Blacklist{}
	return blacklist, db.ODM(ctx).Where(bson.M{"uid": uid, "target_uid": targetUID}).First(blacklist).Error
}

func CreateBlacklist(ctx context.Context, uid, targetUID uint64) (*Blacklist, error) {
	blacklist, err := FindBlacklist(ctx, uid, targetUID)
	if err == nil {
		return blacklist, nil
	}
	if mongo.ErrNoDocuments != err {
		return nil, err
	}
	blacklist = &Blacklist{
		UID:       uid,
		TargetUID: targetUID,
	}
	return blacklist, db.ODM(ctx).Create(blacklist).Error
}

func DeleteBlacklist(ctx context.Context, uid, targetUID uint64) error {
	blacklist, err := FindBlacklist(ctx, uid, targetUID)
	if err != nil {
		return err
	}
	return db.ODM(ctx).Delete(blacklist, blacklist.ID).Error
}
