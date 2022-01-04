package admin

import (
	"context"
	"strconv"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
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
	if params.Tag == enum.TagBlank {
		return a.listAllUser(ctx, params.TraditionalParams)
	}
	tags, page, err := models.ListTag(ctx, &models.ListTagParams{
		PageParams:  params.TraditionalParams,
		TagableType: enum.TagableUser,
		TagType:     params.Tag,
	})
	if err != nil {
		return nil, nil, err
	}
	userIDs := make([]uint64, len(tags))
	for i, tag := range tags {
		uid, _ := strconv.Atoi(tag.TagableID)
		userIDs[i] = uint64(uid)
	}
	users, err := models.ListUserByIDs(ctx, userIDs...)
	if err != nil {
		return nil, nil, err
	}
	userMap := make(map[uint64]*models.User)
	for _, user := range users {
		userMap[user.UID] = user
	}
	result := make([]*UserTag, len(users))
	for i, tag := range tags {
		uid, _ := strconv.Atoi(tag.TagableID)
		result[i] = buildUserTag(userMap[uint64(uid)], tag)
	}
	return result, page, nil
}

func (a *userApi) CreateTag(ctx context.Context, uid uint64, tag enum.TagType) (*UserTag, error) {
	return nil, nil
}

func (a *userApi) DeleteTag(ctx context.Context, uid uint64, tag enum.TagType) (*UserTag, error) {
	return nil, nil
}

func (a *userApi) listAllUser(ctx context.Context, params *pagination.TraditionalParams) ([]*UserTag, pagination.Pagination, error) {
	users := make([]*models.User, 0)
	chain := db.ODM(ctx)
	paginator := pagination.NewTraditionalPaginator(params.PageNum, params.PageSize, chain)
	page, err := paginator.Paginate(&users)
	if err != nil {
		return nil, nil, err
	}
	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = strconv.Itoa(int(user.UID))
	}
	tags, err := models.ListTagByTagables(ctx, enum.TagableUser, userIDs...)
	if err != nil {
		return nil, nil, err
	}
	tagMap := make(map[uint64]*models.Tag)
	for _, tag := range tags {
		uid, _ := strconv.Atoi(tag.TagableID)
		tagMap[uint64(uid)] = tag
	}
	result := make([]*UserTag, len(users))
	for i, user := range users {
		result[i] = buildUserTag(user, tagMap[user.UID])
	}
	return result, page, nil
}

func buildUserTag(user *models.User, tag *models.Tag) *UserTag {
	return &UserTag{
		User: user,
		Tag:  tag,
	}
}
