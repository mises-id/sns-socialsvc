package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	IAdminParams interface {
		BuildAdminSearch(chain *odm.DB) *odm.DB
	}
	IAdminPageParams interface {
		BuildAdminSearch(chain *odm.DB) *odm.DB
		GetPageParams() *pagination.TraditionalParams
	}
	IAdminQuickPageParams interface {
		BuildAdminSearch(chain *odm.DB) *odm.DB
		GetQuickPageParams() *pagination.PageQuickParams
	}
)

func AdminFindStatus(ctx context.Context, params IAdminParams) (*Status, error) {

	status := &Status{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Last(status).Error
	if err != nil {
		return nil, err
	}
	if err := adminHandleStatus(ctx, status); err != nil {
		return nil, err
	}
	return status, preloadStatusUser(ctx, status)
}

func AdminListStatus(ctx context.Context, params IAdminParams) ([]*Status, error) {
	statuses := make([]*Status, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	if err := adminHandleStatus(ctx, statuses...); err != nil {
		return nil, err
	}
	return statuses, preloadStatusUser(ctx, statuses...)
}
func NewListStatus(ctx context.Context, params IAdminParams) ([]*Status, error) {
	statuses := make([]*Status, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, PreloadStatusData(ctx, true, statuses...)
}

func AdminPageStatus(ctx context.Context, params IAdminPageParams) ([]*Status, pagination.Pagination, error) {
	statuses := make([]*Status, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	pageParams := params.GetPageParams()
	paginator := pagination.NewTraditionalPaginatorAdmin(pageParams.PageNum, pageParams.PageSize, chain)
	page, err := paginator.Paginate(&statuses)
	if err != nil {
		return nil, nil, err
	}
	if err := adminHandleStatus(ctx, statuses...); err != nil {
		return nil, nil, err
	}
	return statuses, page, preloadStatusUser(ctx, statuses...)
}

func adminHandleStatus(ctx context.Context, statuses ...*Status) (err error) {

	if err = preloadRelatedStatus(ctx, statuses...); err != nil {
		return err
	}
	if err = preloadAttachment(ctx, statuses...); err != nil {
		return err
	}
	if err = preloadImage(ctx, statuses...); err != nil {
		return err
	}
	return nil
}

func UpdateStatusTag(ctx context.Context, status *Status) error {
	_, err := db.DB().Collection("statuses").UpdateOne(ctx, &bson.M{
		"_id": status.ID,
	}, bson.D{{
		Key: "$set",
		Value: bson.M{
			"tags":       status.Tags,
			"updated_at": time.Now(),
			"score":      time.Now().UnixMilli(),
		}}})
	return err
}

func CountStatus(ctx context.Context, params IAdminParams) (int64, error) {

	var res int64
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Model(&Status{}).Count(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}
