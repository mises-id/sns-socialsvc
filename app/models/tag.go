package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	TagableID   string             `bson:"tagable_id"`
	TagableType enum.TagableType   `bson:"tagable_type"`
	TagType     enum.TagType       `bson:"tag_type"`
	Note        string             `bson:"note"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}

type (
	TagParams struct {
		TagableType enum.TagableType
		TagType     enum.TagType
	}

	ListTagParams struct {
		TagParams
		ListNum int64 `json:"list_num" query:"list_num"`
	}
)
type PageTagParams struct {
	PageParams *pagination.TraditionalParams
	TagParams
}

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

func DeleteTag(ctx context.Context, tagable_id string, tagable_type enum.TagableType, tag_type enum.TagType) error {
	_, err := db.DB().Collection("tags").DeleteOne(ctx, bson.M{"tagable_type": tagable_type, "tag_type": tag_type, "tagable_id": tagable_id})
	return err
}

func DeleteTagsByTagtypes(ctx context.Context, tagable_id string, tagable_type enum.TagableType, tag_types ...enum.TagType) error {
	if tag_types == nil || len(tag_types) == 0 {
		return nil
	}
	_, err := db.DB().Collection("tags").DeleteMany(ctx, bson.M{"tagable_type": tagable_type, "tag_type": bson.M{"$in": tag_types}, "tagable_id": tagable_id})
	return err
}

func ListTagByTagables(ctx context.Context, tagableType enum.TagableType, tagableIDs ...string) ([]*Tag, error) {
	tags := make([]*Tag, 0)
	return tags, db.ODM(ctx).Where(bson.M{"tagable_type": tagableType, "tagable_id": bson.M{"$in": tagableIDs}}).Find(&tags).Error
}

func ListTag(ctx context.Context, params *ListTagParams) ([]*Tag, error) {
	tags := make([]*Tag, 0)
	return tags, db.ODM(ctx).Where(bson.M{"tagable_type": params.TagableType, "tag_type": params.TagType}).Sort(bson.M{"_id": -1}).Find(&tags).Error
}

func PageTag(ctx context.Context, params *PageTagParams) ([]*Tag, pagination.Pagination, error) {
	tags := make([]*Tag, 0)
	chain := db.ODM(ctx).Where(bson.M{"tagable_type": params.TagableType, "tag_type": params.TagType}).Sort(bson.M{"_id": -1})
	paginator := pagination.NewTraditionalPaginator(params.PageParams.PageNum, params.PageParams.PageSize, chain)
	page, err := paginator.Paginate(&tags)
	if err != nil {
		return nil, nil, err
	}
	return tags, page, nil
}
