package models

import (
	"context"
	"time"

	"github.com/mises-id/socialsvc/lib/db"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func EnsureIndex() {
	opts := options.CreateIndexes().SetMaxTime(20 * time.Second)
	trueBool := true
	_, err := db.DB().Collection("users").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{"username": 1},
			Options: &options.IndexOptions{
				Unique: &trueBool,
				Sparse: &trueBool,
			},
		},
		{
			Keys: bson.M{"misesid": 1},
			Options: &options.IndexOptions{
				Unique: &trueBool,
			},
		},
	}, opts)
	if err != nil {
		logrus.Debug(err)
	}
	_, err = db.DB().Collection("follows").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{"to_uid": 1},
		},
		{
			Keys: bsonx.Doc{{
				Key: "from_uid", Value: bsonx.Int32(1),
			}, {
				Key: "to_uid", Value: bsonx.Int32(1)},
			},
			Options: &options.IndexOptions{
				Unique: &trueBool,
			},
		},
	}, opts)
	if err != nil {
		logrus.Debug(err)
	}

	_, err = db.DB().Collection("statuses").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.M{"uid": 1},
		},
		{
			Keys: bson.M{"deleted_at": 1},
		},
	}, opts)
	if err != nil {
		logrus.Debug(err)
	}

	_, err = db.DB().Collection("likes").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{
				Key: "uid", Value: bsonx.Int32(1),
			}, {
				Key: "target_id", Value: bsonx.Int32(1),
			}, {
				Key: "target_type", Value: bsonx.Int32(1),
			}, {
				Key: "deleted_at", Value: bsonx.Int32(1)},
			},
			Options: &options.IndexOptions{
				Unique: &trueBool,
			},
		},
	}, opts)
	if err != nil {
		logrus.Debug(err)
	}
}
