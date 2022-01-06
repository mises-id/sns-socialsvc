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
	ListStatusTagParams struct {
		PageNum   int64      `json:"page_num" query:"page_num"`
		PageSize  int64      `json:"page_size" query:"page_size"`
		StartTime *time.Time `json:"start_time" query:"start_time"` //发布起始时间
		EndTime   *time.Time `json:"end_time" query:"end_time"`
		Tag       enum.TagType
	}
	AdminStatusParams struct {
		//search
		IDs       []primitive.ObjectID
		UIDs      []uint64
		FromTypes []enum.FromType
		Tags      []enum.TagType
		StartTime *time.Time `json:"start_time" query:"start_time"`
		EndTime   *time.Time `json:"end_time" query:"end_time"`
		Tag       enum.TagType
		//sort
		//limit
		ListNum int64
		//page
		PageNum  int64 `json:"page_num" query:"page_num"`
		PageSize int64 `json:"page_size" query:"page_size"`
		//PageParams *pagination.TraditionalParams
	}
)

func (params *AdminStatusParams) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//where
	if params.IDs != nil && len(params.IDs) > 0 {
		chain = chain.Where(bson.M{"_id": bson.M{"$in": params.IDs}})
	}
	if params.UIDs != nil && len(params.UIDs) > 0 {
		chain = chain.Where(bson.M{"uid": bson.M{"$in": params.UIDs}})
	}
	if params.FromTypes != nil {
		chain = chain.Where(bson.M{"from_type": bson.M{"$in": params.FromTypes}})
	}
	if params.Tag != enum.TagBlank {
		params.Tags = []enum.TagType{params.Tag}
	}
	if params.Tags != nil {
		chain = chain.Where(bson.M{"tags": bson.M{"$in": params.Tags}})
	}
	if params.StartTime != nil {
		chain = chain.Where(bson.M{"create_at": bson.M{"$gte": params.StartTime}})
	}
	if params.EndTime != nil {
		chain = chain.Where(bson.M{"create_at": bson.M{"$lte": params.EndTime}})
	}
	//sort
	//limit
	if (params.PageNum <= 0 || params.PageSize <= 0) && params.ListNum > 0 {
		chain = chain.Limit(params.ListNum)
	}
	return chain
}

func (params *AdminStatusParams) GetPageParams() *pagination.TraditionalParams {
	if params.PageNum <= 0 || params.PageSize <= 0 {
		return pagination.DefaultTraditionalParams()
	}
	return &pagination.TraditionalParams{
		PageNum:  params.PageNum,
		PageSize: params.PageSize,
	}
}
