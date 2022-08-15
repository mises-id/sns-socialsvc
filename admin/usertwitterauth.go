package admin

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
)

type (
	UserTwitterAuth struct {
		*models.UserTwitterAuth
	}
)
type UserTwitterAuthApi interface {
	PageUserTwitterAuth(ctx context.Context, params *search.UserTwitterAuthSearch) ([]*UserTwitterAuth, pagination.Pagination, error)
}

type userTwitterAuthApi struct {
}

func NewUserTwitterAuthApi() UserTwitterAuthApi {
	return &userTwitterAuthApi{}
}

//get page
func (a *userTwitterAuthApi) PageUserTwitterAuth(ctx context.Context, params *search.UserTwitterAuthSearch) ([]*UserTwitterAuth, pagination.Pagination, error) {

	//var res []*models.UserTwitterAuth
	var page pagination.Pagination
	//var err error
	return nil, page, nil
	/* res, page, err = models.AdminPageUserTwitterAuth(ctx, params)
	if err != nil {
		return nil, nil, err
	}
	result := make([]*UserTwitterAuth, len(res))
	for i, v := range res {
		result[i] = &UserTwitterAuth{v}
	}
	return result, page, nil */
}
