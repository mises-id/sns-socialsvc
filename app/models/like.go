package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty"`
	UID        uint64              `bson:"uid,omitempty"`
	TargetID   primitive.ObjectID  `bson:"target_id,omitempty"`
	TargetType enum.LikeTargetType `bson:"target_type"`
	DeletedAt  time.Time           `bson:"deleted_at,omitempty"`
	CreatedAt  time.Time           `bson:"created_at,omitempty"`
	UpdatedAt  time.Time           `bson:"updated_at,omitempty"`
}

func CreateLike(ctx context.Context, uid uint64, targetID primitive.ObjectID, targetType enum.LikeTargetType) (*Like, error) {
	like := &Like{
		UID:        uid,
		TargetID:   targetID,
		TargetType: targetType,
	}
	return like, db.ODM(ctx).Create(like).Error
}

func DeleteLike(ctx context.Context, id primitive.ObjectID) error {
	return db.DB().Collection("likes").FindOneAndUpdate(ctx, bson.M{
		"_id":        id,
		"deleted_at": bson.M{"$exists": false},
	}, bson.M{"$set": bson.M{"deleted_at": time.Now()}}).Err()
}

func FindLike(ctx context.Context, uid uint64, targetID primitive.ObjectID, targetType enum.LikeTargetType) (*Like, error) {
	like := &Like{}
	err := db.ODM(ctx).Where(bson.M{
		"uid":         uid,
		"target_id":   targetID,
		"target_type": targetType,
		"deleted_at":  nil,
	}).First(like).Error
	return like, err
}

func GetStatusLikeMap(ctx context.Context, uid uint64, statusIDs []primitive.ObjectID) (map[primitive.ObjectID]*Like, error) {
	likes := make([]*Like, 0)
	err := db.ODM(ctx).Where(bson.M{
		"uid":         uid,
		"target_id":   bson.M{"$in": statusIDs},
		"target_type": enum.LikeStatus,
		"deleted_at":  nil,
	}).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	likeMap := make(map[primitive.ObjectID]*Like)
	for _, like := range likes {
		likeMap[like.TargetID] = like
	}
	return likeMap, nil
}
