package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	AdminPageUserParams struct {
		IDs        []uint64
		UIDs       []uint64
		Tags       []enum.TagType
		PageParams *pagination.TraditionalParams
	}
)

func AdminFindUser(ctx context.Context, params IAdminParams) (*User, error) {

	user := &User{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Last(user).Error
	if err != nil {
		return nil, err
	}

	return user, PreloadUserAvatar(ctx, user)
}

func AdminListUser(ctx context.Context, params IAdminParams) ([]*User, error) {

	users := make([]*User, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, PreloadUserAvatar(ctx, users...)
}

func AdminPageUser(ctx context.Context, params IAdminPageParams) ([]*User, pagination.Pagination, error) {
	users := make([]*User, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	pageParams := params.GetPageParams()
	paginator := pagination.NewTraditionalPaginator(pageParams.PageNum, pageParams.PageSize, chain)
	page, err := paginator.Paginate(&users)
	if err != nil {
		return nil, nil, err
	}

	return users, page, PreloadUserAvatar(ctx, users...)
}

func UpdateUserTag(ctx context.Context, user *User) error {
	_, err := db.DB().Collection("users").UpdateOne(ctx, &bson.M{
		"_id": user.UID,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"tags":       user.Tags,
			"updated_at": time.Now(),
		}}})
	return err
}
