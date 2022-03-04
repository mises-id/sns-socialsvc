package models

import (
	"context"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Config struct {
		ID    primitive.ObjectID `bson:"_id,omitempty"`
		Key   string             `bson:"key"`
		Value interface{}        `bson:"value"`
	}
)

func FindOrCreateConfig(ctx context.Context, key string, value interface{}) (*Config, error) {
	res, err := findConfigByKey(ctx, key)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return createdConfig(ctx, key, value)
		}
		return nil, err
	}
	return res, nil
}
func UpdateOrCreateConfig(ctx context.Context, key string, value interface{}) error {
	res, err := FindOrCreateConfig(ctx, key, value)
	if err != nil {
		return err
	}
	return UpdateConfig(ctx, res.ID, value)
}

func UpdateConfig(ctx context.Context, id primitive.ObjectID, value interface{}) error {
	_, err := db.DB().Collection("configs").UpdateOne(ctx, &bson.M{
		"_id": id,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"value": value,
		}}})
	return err
}

func createdConfig(ctx context.Context, key string, value interface{}) (*Config, error) {
	created := &Config{
		Key: key, Value: value,
	}
	res, err := db.DB().Collection("configs").InsertOne(ctx, created)
	if err != nil {
		return nil, err
	}
	created.ID = res.InsertedID.(primitive.ObjectID)
	return created, err
}

func findConfigByKey(ctx context.Context, key string) (*Config, error) {
	res := &Config{}
	chain := db.ODM(ctx).Where(bson.M{"key": key})
	err := chain.Last(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
