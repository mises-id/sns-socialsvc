package search

import (
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	NftLogSearch struct {
		ID             primitive.ObjectID
		NftTagableType enum.NftTagableType
		ObjectID       string
		LastID         primitive.ObjectID
		//sort
		SortBy string
		//limit
		ListNum int64
		//page
		PageNum    int64 `json:"page_num" query:"page_num"`
		PageSize   int64 `json:"page_size" query:"page_size"`
		PageParams *pagination.PageQuickParams
	}
)

func (params *NftLogSearch) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//base
	//where

	if !params.ID.IsZero() {
		chain = chain.Where(bson.M{"_id": params.ID})
	}
	if params.NftTagableType != enum.NftTagableTypeEmpty {
		chain = chain.Where(bson.M{"nft_tagable_type": params.NftTagableType})
	}
	if params.ObjectID != "" {
		chain = chain.Where(bson.M{"object_id": params.ObjectID})
	}
	if !params.LastID.IsZero() {
		chain = chain.Where(bson.M{"_id": bson.M{"$lte": params.LastID}})
	}
	//sort

	chain = chain.Sort(bson.M{"_id": -1})
	//limit
	if (params.PageNum <= 0 || params.PageSize <= 0) && params.ListNum > 0 {
		chain = chain.Limit(params.ListNum)
	}
	return chain
}

func (params *NftLogSearch) GetPageParams() *pagination.TraditionalParams {
	if params.PageNum <= 0 || params.PageSize <= 0 {
		return pagination.DefaultTraditionalParams()
	}
	return &pagination.TraditionalParams{
		PageNum:  params.PageNum,
		PageSize: params.PageSize,
	}
}
func (params *NftLogSearch) GetQuickPageParams() *pagination.PageQuickParams {
	if params.PageParams == nil {
		return pagination.DefaultQuickParams()
	}
	return params.PageParams
}
