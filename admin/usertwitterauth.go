package admin

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	UserTwitterAuth struct {
		*models.UserTwitterAuth
	}
	RankUserTwitterAuthParams struct {
		Pipe bson.A
	}
	UserTwitterAuthRank struct {
		*models.UserTwitterAuthRank
	}
)
type UserTwitterAuthApi interface {
	PageUserTwitterAuth(ctx context.Context, params *search.UserTwitterAuthSearch) ([]*UserTwitterAuth, pagination.Pagination, error)
	RankUserTwitterAuth(ctx context.Context, params *RankUserTwitterAuthParams) ([]*UserTwitterAuthRank, error)
	Count(ctx context.Context, params *search.UserTwitterAuthSearch) (int64, error)
}

type userTwitterAuthApi struct {
}

func NewUserTwitterAuthApi() UserTwitterAuthApi {
	return &userTwitterAuthApi{}
}

//get page
func (a *userTwitterAuthApi) PageUserTwitterAuth(ctx context.Context, params *search.UserTwitterAuthSearch) ([]*UserTwitterAuth, pagination.Pagination, error) {

	var res []*models.UserTwitterAuth
	var page pagination.Pagination
	var err error

	res, page, err = models.AdminPageUserTwitterAuth(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	result := make([]*UserTwitterAuth, len(res))
	for i, v := range res {
		result[i] = &UserTwitterAuth{v}
	}
	return result, page, nil
}

func (a *userTwitterAuthApi) RankUserTwitterAuth(ctx context.Context, params *RankUserTwitterAuthParams) ([]*UserTwitterAuthRank, error) {

	ranks, err := models.RankUserTwitterAuth(ctx, &models.RankUserTwitterAuthParams{Pipe: params.Pipe})
	if err != nil {
		return nil, err
	}
	result := make([]*UserTwitterAuthRank, len(ranks))
	for i, rank := range ranks {
		result[i] = &UserTwitterAuthRank{rank}
	}
	return result, nil
}

func (api *userTwitterAuthApi) Count(ctx context.Context, params *search.UserTwitterAuthSearch) (int64, error) {
	return models.CountUserTwitterAuth(ctx, params)
}
