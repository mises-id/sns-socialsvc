package models

import (
	"context"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AdminListBlackListUserIDs(ctx context.Context, uid uint64) ([]uint64, error) {
	cursor, err := db.DB().Collection("blacklists").Find(ctx, &bson.M{
		"uid": uid,
	}, &options.FindOptions{
		Projection: bson.M{"target_uid": 1},
	})
	if err != nil {
		return nil, err
	}
	blacklists := make([]*Blacklist, 0)
	if err = cursor.All(ctx, &blacklists); err != nil {
		return nil, err
	}
	ids := make([]uint64, len(blacklists))
	for i, blacklist := range blacklists {
		ids[i] = blacklist.TargetUID
	}
	return ids, nil
}
