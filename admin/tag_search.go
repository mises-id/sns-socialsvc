package admin

import (
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AdminTagParams struct {
		//search
		ID          primitive.ObjectID
		IDs         []primitive.ObjectID
		TagableIDs  []string
		TagableType enum.TagableType
		TagTypes    []enum.TagType
		StartTime   *time.Time `json:"start_time" query:"start_time"`
		EndTime     *time.Time `json:"end_time" query:"end_time"`
		//sort
		SortKey  string `json:"sort_key" query:"sort_key" validate:"omitempty,oneof=created_at comments_count likes_count forwards_count"` //发布时间/评论数/点赞数/转发数
		SortType int    `json:"sort_type" query:"sort_type" validate:"omitempty,oneof=-1 1"`
		//limit
		ListNum int64
		//page
		PageNum  int64 `json:"page_num" query:"page_num"`
		PageSize int64 `json:"page_size" query:"page_size"`
		//PageParams *pagination.TraditionalParams
	}
)

func (params *AdminTagParams) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//base
	chain = chain.Sort(bson.M{"_id": -1})
	//where
	if !params.ID.IsZero() {
		params.IDs = []primitive.ObjectID{params.ID}
	}
	if params.IDs != nil && len(params.IDs) > 0 {
		chain = chain.Where(bson.M{"_id": bson.M{"$in": params.IDs}})
	}
	if params.TagableIDs != nil && len(params.TagableIDs) > 0 {
		chain = chain.Where(bson.M{"tagable_id": bson.M{"$in": params.TagableIDs}})
	}
	if params.TagTypes != nil && len(params.TagTypes) > 0 {
		chain = chain.Where(bson.M{"tag_type": bson.M{"$in": params.TagTypes}})
	}
	if params.TagableType != "" {
		chain = chain.Where(bson.M{"tagable_type": params.TagableType})
	}
	if params.StartTime != nil && params.EndTime == nil {
		chain = chain.Where(bson.M{"created_at": bson.M{"$gte": params.StartTime}})
	}
	if params.StartTime == nil && params.EndTime != nil {
		chain = chain.Where(bson.M{"created_at": bson.M{"$lte": params.EndTime}})
	}
	if params.StartTime != nil && params.EndTime != nil {
		chain = chain.Where(bson.M{"$and": bson.A{bson.M{"created_at": bson.M{"$gte": params.StartTime}}, bson.M{"created_at": bson.M{"$lte": params.EndTime}}}})
	}
	//sort
	if params.SortKey != "" && params.SortType != 0 {
		chain = chain.Sort(bson.M{params.SortKey: params.SortType})
	}
	//limit
	if (params.PageNum <= 0 || params.PageSize <= 0) && params.ListNum > 0 {
		chain = chain.Limit(params.ListNum)
	}
	return chain
}

func (params *AdminTagParams) GetPageParams() *pagination.TraditionalParams {
	if params.PageNum <= 0 || params.PageSize <= 0 {
		return pagination.DefaultTraditionalParams()
	}
	return &pagination.TraditionalParams{
		PageNum:  params.PageNum,
		PageSize: params.PageSize,
	}
}
