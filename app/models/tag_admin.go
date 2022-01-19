package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *Tag) BeforeCreate(ctx context.Context) error {
	var err error
	if err != nil {
		return err
	}
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return m.Validate(ctx)
}

func (m *Tag) Validate(ctx context.Context) error {
	if m.TagType == "" {
		return codes.ErrInvalidArgument
	}
	if m.TagableID == "" {
		return codes.ErrInvalidArgument
	}
	if m.TagableType == "" {
		return codes.ErrInvalidArgument
	}
	var n int64
	if err := db.ODM(ctx).Model(m).Where(bson.M{"tagable_type": m.TagableType, "tag_type": m.TagType, "tagable_id": m.TagableID}).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return codes.ErrInvalidArgument.New("record exists")
	}
	return nil
}

type CreateTagParams struct {
	TagableID   string           `bson:"tagable_id"`
	TagableType enum.TagableType `bson:"tagable_type"`
	TagType     enum.TagType     `bson:"tag_type"`
	Note        string           `bson:"note"`
}

func CreateTag(ctx context.Context, params *CreateTagParams) (*Tag, error) {
	data := &Tag{
		TagableType: params.TagableType,
		TagableID:   params.TagableID,
		TagType:     params.TagType,
		Note:        params.Note,
	}
	var err error
	if err = data.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	if err = db.ODM(ctx).Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func AdminFindTag(ctx context.Context, params IAdminParams) (*Tag, error) {

	tag := &Tag{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Last(tag).Error
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func AdminListTag(ctx context.Context, params IAdminParams) ([]*Tag, error) {

	tags := make([]*Tag, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func AdminPageTag(ctx context.Context, params IAdminPageParams) ([]*Tag, pagination.Pagination, error) {
	tags := make([]*Tag, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	pageParams := params.GetPageParams()
	paginator := pagination.NewTraditionalPaginatorAdmin(pageParams.PageNum, pageParams.PageSize, chain)
	page, err := paginator.Paginate(&tags)
	if err != nil {
		return nil, nil, err
	}

	return tags, page, nil
}

func DeleteTag(ctx context.Context, tagable_id string, tagable_type enum.TagableType, tag_type enum.TagType) error {
	_, err := db.DB().Collection("tags").DeleteOne(ctx, bson.M{"tagable_type": tagable_type, "tag_type": tag_type, "tagable_id": tagable_id})
	return err
}

func DeleteTagsByTagtypes(ctx context.Context, tagable_id string, tagableType enum.TagableType, tagTypes ...enum.TagType) error {
	if len(tagTypes) == 0 {
		return nil
	}
	_, err := db.DB().Collection("tags").DeleteMany(ctx, bson.M{"tagable_type": tagableType, "tag_type": bson.M{"$in": tagTypes}, "tagable_id": tagable_id})
	return err
}
