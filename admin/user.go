package admin

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
)

type ListUserTagParams struct {
	*pagination.TraditionalParams
	Tag enum.TagType
}

type UserTag struct {
	*models.User
	*models.Tag
}

type UserApi interface {
	ListTag(ctx context.Context, params *ListUserTagParams) ([]*UserTag, pagination.Pagination, error)
	CreateTag(ctx context.Context, uid uint64, tag enum.TagType) (*UserTag, error)
	DeleteTag(ctx context.Context, uid uint64, tag enum.TagType) (*UserTag, error)
}

type userApi struct {
}

func NewUserApi() UserApi {
	return &userApi{}
}

func (a *userApi) ListTag(ctx context.Context, params *ListUserTagParams) ([]*UserTag, pagination.Pagination, error) {
	return nil, nil, nil
}

func (a *userApi) CreateTag(ctx context.Context, uid uint64, tag enum.TagType) (*UserTag, error) {
	return nil, nil
}

func (a *userApi) DeleteTag(ctx context.Context, uid uint64, tag enum.TagType) (*UserTag, error) {
	return nil, nil
}
