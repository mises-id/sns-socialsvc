package admin

import (
	"fmt"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AdminStatusParams struct {
		//search
		ID        primitive.ObjectID
		IDs       []primitive.ObjectID
		UIDs      []uint64
		FromTypes []enum.FromType
		Tags      []enum.TagType
		StartTime *time.Time `json:"start_time" query:"start_time"`
		EndTime   *time.Time `json:"end_time" query:"end_time"`
		Tag       enum.TagType
		OnlyShow  bool
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

func (params *AdminStatusParams) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//base
	chain = chain.Sort(bson.M{"_id": -1})
	//where
	if !params.ID.IsZero() {
		params.IDs = []primitive.ObjectID{params.ID}
	}
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
	if params.StartTime != nil && params.EndTime == nil {
		fmt.Println("st:", params.StartTime)
		chain = chain.Where(bson.M{"created_at": bson.M{"$gte": params.StartTime}})
	}
	if params.StartTime == nil && params.EndTime != nil {
		fmt.Println("et:", params.EndTime)
		chain = chain.Where(bson.M{"created_at": bson.M{"$lte": params.EndTime}})
	}
	if params.StartTime != nil && params.EndTime != nil {
		fmt.Println("et:", params.EndTime)
		chain = chain.Where(bson.M{"$and": bson.A{bson.M{"created_at": bson.M{"$gte": params.StartTime}}, bson.M{"created_at": bson.M{"$lte": params.EndTime}}}})
	}
	if params.OnlyShow {
		chain = chain.Where(bson.M{"$or": bson.A{bson.M{"hide_time": nil}, bson.M{"hide_time": bson.M{"$gt": time.Now()}}}})
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

func (params *AdminStatusParams) GetPageParams() *pagination.TraditionalParams {
	if params.PageNum <= 0 || params.PageSize <= 0 {
		return pagination.DefaultTraditionalParams()
	}
	return &pagination.TraditionalParams{
		PageNum:  params.PageNum,
		PageSize: params.PageSize,
	}
}
