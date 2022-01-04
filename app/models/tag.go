package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
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

type ListTagParams struct {
	PageParams  *pagination.TraditionalParams
	TagableType enum.TagableType
	TagType     enum.TagType
}

func ListTagByTagables(ctx context.Context, tagableType enum.TagableType, tagableIDs ...string) ([]*Tag, error) {
	tags := make([]*Tag, 0)
	return tags, db.ODM(ctx).Where(bson.M{"tagable_type": tagableType, "tagable_id": bson.M{"$in": tagableIDs}}).Find(&tags).Error
}

func ListTag(ctx context.Context, params *ListTagParams) ([]*Tag, pagination.Pagination, error) {
	tags := make([]*Tag, 0)
	chain := db.ODM(ctx).Where(bson.M{"tagable_type": params.TagableType, "tag_type": params.TagType}).Sort(bson.M{"_id": -1})
	paginator := pagination.NewTraditionalPaginator(params.PageParams.PageNum, params.PageParams.PageSize, chain)
	page, err := paginator.Paginate(&tags)
	if err != nil {
		return nil, nil, err
	}
	return tags, page, nil
}
