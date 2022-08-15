package admin

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	AirdropApi interface {
		RankAirdrop(ctx context.Context, params *RankAirdropParams) ([]*AirdropRank, error)
		Count(ctx context.Context, params *AirdropSearch) (int64, error)
	}

	RankAirdropParams struct {
		Pipe bson.A
	}
	AirdropRank struct {
		*models.AirdropRank
	}
	airdropApi struct{}
)

func NewAirdropApi() AirdropApi {
	return &airdropApi{}
}

func (a *airdropApi) RankAirdrop(ctx context.Context, params *RankAirdropParams) ([]*AirdropRank, error) {

	ranks, err := models.RankAirdrop(ctx, &models.RankAirdropParams{Pipe: params.Pipe})
	if err != nil {
		return nil, err
	}
	result := make([]*AirdropRank, len(ranks))
	for i, rank := range ranks {
		result[i] = &AirdropRank{rank}
	}
	return result, nil
}

func (api *airdropApi) Count(ctx context.Context, params *AirdropSearch) (int64, error) {
	return models.CountAirdrop(ctx, params)
}
