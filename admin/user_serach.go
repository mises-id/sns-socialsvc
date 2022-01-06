package admin

import (
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	AdminUserParams struct {
		//search
		IDs       []uint64
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

func (params *AdminUserParams) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//where
	if params.IDs != nil && len(params.IDs) > 0 {
		chain = chain.Where(bson.M{"_id": bson.M{"$in": params.IDs}})
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

func (params *AdminUserParams) GetPageParams() *pagination.TraditionalParams {
	if params.PageNum <= 0 || params.PageSize <= 0 {
		return pagination.DefaultTraditionalParams()
	}
	return &pagination.TraditionalParams{
		PageNum:  params.PageNum,
		PageSize: params.PageSize,
	}
}
