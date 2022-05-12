package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	AdminPageUserParams struct {
		IDs        []uint64
		UIDs       []uint64
		Tags       []enum.TagType
		PageParams *pagination.TraditionalParams
	}
)

//find one user
func AdminFindUser(ctx context.Context, params IAdminParams) (*User, error) {

	user := &User{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Last(user).Error
	if err != nil {
		return nil, err
	}

	return user, preloadUserAvatar(ctx, user)
}

//list user
func AdminListUser(ctx context.Context, params IAdminParams) ([]*User, error) {

	users := make([]*User, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, preloadUserAvatar(ctx, users...)
}

//page user
func AdminPageUser(ctx context.Context, params IAdminPageParams) ([]*User, pagination.Pagination, error) {
	users := make([]*User, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	pageParams := params.GetPageParams()
	paginator := pagination.NewTraditionalPaginatorAdmin(pageParams.PageNum, pageParams.PageSize, chain)
	page, err := paginator.Paginate(&users)
	if err != nil {
		return nil, nil, err
	}

	return users, page, preloadUserAvatar(ctx, users...)
}

//find problem user uids
func AdminListProblemUserIDs(ctx context.Context) ([]uint64, error) {
	cursor, err := db.DB().Collection("users").Find(ctx, &bson.M{
		"tags": bson.M{"$in": []enum.TagType{enum.TagProblemUser}},
	}, &options.FindOptions{
		Projection: bson.M{"id": 1},
	})
	if err != nil {
		return nil, err
	}
	users := make([]*User, 0)
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	ids := make([]uint64, len(users))
	for i, user := range users {
		ids[i] = user.UID
	}
	return ids, nil
}

//find star user uids
func AdminListStarUserIDs(ctx context.Context) ([]uint64, error) {
	cursor, err := db.DB().Collection("users").Find(ctx, &bson.M{
		"tags": bson.M{"$in": []enum.TagType{enum.TagStarUser}},
	}, &options.FindOptions{
		Projection: bson.M{"id": 1},
	})
	if err != nil {
		return nil, err
	}
	users := make([]*User, 0)
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	ids := make([]uint64, len(users))
	for i, user := range users {
		ids[i] = user.UID
	}
	return ids, nil
}

//update user tags
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

func CountUser(ctx context.Context, params IAdminParams) (int64, error) {

	var res int64
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Model(&User{}).Count(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}
