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
	AdminUserParams struct {
		//search
		ID            uint64
		IDs           []uint64
		FromTypes     []enum.FromType
		Username      string
		MisesID       string
		ChannelID     primitive.ObjectID
		ChannelIDs    []primitive.ObjectID
		IsChannelUser uint64
		AvatarType    uint64
		Tags          []enum.TagType
		StartTime     *time.Time `json:"start_time" query:"start_time"`
		EndTime       *time.Time `json:"end_time" query:"end_time"`
		Tag           enum.TagType
		IsNftUser     bool
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
	//base

	//where
	if params.ID != 0 {
		params.IDs = []uint64{params.ID}
	}
	if params.IDs != nil && len(params.IDs) > 0 {
		chain = chain.Where(bson.M{"_id": bson.M{"$in": params.IDs}})
	}
	if !params.ChannelID.IsZero() {
		chain = chain.Where(bson.M{"channel_id": params.ChannelID})
	}
	if params.ChannelIDs != nil && len(params.ChannelIDs) > 0 {
		chain = chain.Where(bson.M{"channel_id": bson.M{"$in": params.ChannelIDs}})
	}
	if params.AvatarType == 1 {
		chain = chain.Where(bson.M{"nft_avatar": bson.M{"$eq": nil}})
	}
	if params.AvatarType == 2 {
		chain = chain.Where(bson.M{"nft_avatar": bson.M{"$ne": nil}})
	}
	if params.IsChannelUser == 1 {
		chain = chain.Where(bson.M{"channel_id": bson.M{"$ne": nil}})
	}
	if params.IsChannelUser == 2 {
		chain = chain.Where(bson.M{"channel_id": nil})
	}
	if params.Username != "" {
		chain = chain.Where(bson.M{"username": params.Username})
	}
	if params.MisesID != "" {
		chain = chain.Where(bson.M{"misesid": params.MisesID})
	}
	if params.Tag != enum.TagBlank {
		params.Tags = []enum.TagType{params.Tag}
	}
	if params.Tags != nil {
		chain = chain.Where(bson.M{"tags": bson.M{"$in": params.Tags}})
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
	chain = chain.Sort(bson.M{"_id": -1})
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
