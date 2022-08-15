package models

import (
	"context"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	AirdropRank struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
		Coin  uint64 `bson:"coin"`
	}
	RankAirdropParams struct {
		Pipe bson.A
	}
)

func RankAirdrop(ctx context.Context, params *RankAirdropParams) ([]*AirdropRank, error) {
	out := make([]*AirdropRank, 0)
	pipe := bson.A{
		bson.M{"$match": bson.M{"status": 2}},
		bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
	}
	if params.Pipe != nil && len(params.Pipe) > 0 {
		pipe = params.Pipe
	}
	res, err := db.DB().Collection("airdrops").Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}
	err = res.All(ctx, &out)
	return out, err
}

func CountAirdrop(ctx context.Context, params IAdminParams) (int64, error) {

	var res int64
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Model(&Airdrop{}).Count(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}
