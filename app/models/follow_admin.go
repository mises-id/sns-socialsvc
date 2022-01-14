package models

import (
	"context"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ()

func AdminListFollowingUserIDs(ctx context.Context, uids []uint64) ([]uint64, error) {
	cursor, err := db.DB().Collection("follows").Find(ctx, &bson.M{
		"from_uid": bson.M{"$in": uids},
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
