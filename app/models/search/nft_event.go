package search

import (
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	NftEventSearch struct {
		ID         primitive.ObjectID
		NftAssetID primitive.ObjectID
		LastID     primitive.ObjectID
		EventTypes []string
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

func (params *NftEventSearch) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//base
	//where

	if !params.ID.IsZero() {
		chain = chain.Where(bson.M{"_id": params.ID})
	}
	if !params.NftAssetID.IsZero() {
		chain = chain.Where(bson.M{"nft_asset_id": params.NftAssetID})
	}

	if !params.LastID.IsZero() {
		chain = chain.Where(bson.M{"_id": bson.M{"$lte": params.LastID}})
	}
	if params.EventTypes != nil && len(params.EventTypes) > 0 {
		chain = chain.Where(bson.M{"event_type": bson.M{"$in": params.EventTypes}})
	}
	//sort
	if params.SortBy != "" {
		switch params.SortBy {
		case "created_date_desc":
			chain = chain.Sort(bson.M{"created_date": -1})
		}

	}
	chain = chain.Sort(bson.M{"_id": -1})
	//limit
	if (params.PageNum <= 0 || params.PageSize <= 0) && params.ListNum > 0 {
		chain = chain.Limit(params.ListNum)
	}
	return chain
}

func (params *NftEventSearch) GetPageParams() *pagination.TraditionalParams {
	if params.PageNum <= 0 || params.PageSize <= 0 {
		return pagination.DefaultTraditionalParams()
	}
	return &pagination.TraditionalParams{
		PageNum:  params.PageNum,
		PageSize: params.PageSize,
	}
}
func (params *NftEventSearch) GetQuickPageParams() *pagination.PageQuickParams {
	if params.PageParams == nil {
		return pagination.DefaultQuickParams()
	}
	return params.PageParams
}
