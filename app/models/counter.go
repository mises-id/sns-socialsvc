package models

import (
	"context"
	"sync"

	"github.com/mises-id/socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var mtx sync.Mutex

type Counter struct {
	ID  string `bson:"_id,omitempty"`
	Seq uint64 `bson:"seq,omitempty"`
}

func getNextSeq(ctx context.Context, id string) (uint64, error) {
	seq, err := incSeq(ctx, id)
	if err == nil {
		return seq, nil
	}
	if mongo.ErrNoDocuments == err {
		return firstSeq(ctx, id)
	}
	return 0, err
}

func incSeq(ctx context.Context, id string) (uint64, error) {
	result := db.DB().Collection("counters").FindOneAndUpdate(ctx, bson.M{"_id": id},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   "seq",
				Value: 1,
			}}},
		})
	err := result.Err()
	if err != nil {
		return 0, err
	}
	counter := &Counter{}
	return counter.Seq + 1, result.Decode(counter)
}

func firstSeq(ctx context.Context, id string) (uint64, error) {
	mtx.Lock()
	defer mtx.Unlock()
	seq, err := incSeq(ctx, id)
	if err == nil {
		return seq, nil
	}
	if mongo.ErrNoDocuments != err {
		return 0, err
	}
	counter := &Counter{
		ID:  id,
		Seq: 1,
	}
	_, err = db.DB().Collection("counters").InsertOne(ctx, counter)
	return counter.Seq, err
}
