package admin

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	ChannelUserApi interface {
		RankChannelUser(ctx context.Context, params *RankChannelUserParams) ([]*ChannelUserRank, error)
		Count(ctx context.Context, params *ChannelUserSearch) (int64, error)
	}

	RankChannelUserParams struct {
		Pipe bson.A
	}
	ChannelUserRank struct {
		*models.ChannelUserRank
	}
	channelUserApi struct{}
)

func NewChannelUserApi() ChannelUserApi {
	return &channelUserApi{}
}

func (a *channelUserApi) RankChannelUser(ctx context.Context, params *RankChannelUserParams) ([]*ChannelUserRank, error) {

	ranks, err := models.RankChannelUser(ctx, &models.RankChannelUserParams{Pipe: params.Pipe})
	if err != nil {
		return nil, err
	}
	result := make([]*ChannelUserRank, len(ranks))
	for i, rank := range ranks {
		result[i] = &ChannelUserRank{rank}
	}
	return result, nil
}

func (api *channelUserApi) Count(ctx context.Context, params *ChannelUserSearch) (int64, error) {
	return models.CountChannelUser(ctx, params)
}
