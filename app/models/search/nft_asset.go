package search

import (
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	NftAssetSearch struct {
		ID     primitive.ObjectID
		UID    uint64
		LastID primitive.ObjectID
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

func (params *NftAssetSearch) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//base
	//where

	if !params.ID.IsZero() {
		chain = chain.Where(bson.M{"_id": params.ID})
	}
	if params.UID > 0 {
		chain = chain.Where(bson.M{"uid": params.UID})
	}
	if !params.LastID.IsZero() {
		chain = chain.Where(bson.M{"_id": bson.M{"$lte": params.LastID}})
	}
	//sort
	if params.SortBy != "" {
		switch params.SortBy {
		case "collection_asc":
			chain = chain.Sort(bson.M{"collection.slug": 1})
		case "collection_desc":
			chain = chain.Sort(bson.M{"collection.slug": -1})
		}

	}
	chain = chain.Sort(bson.M{"_id": -1})
	//limit
	if (params.PageNum <= 0 || params.PageSize <= 0) && params.ListNum > 0 {
		chain = chain.Limit(params.ListNum)
	}
	return chain
}

func (params *NftAssetSearch) GetPageParams() *pagination.TraditionalParams {
	if params.PageNum <= 0 || params.PageSize <= 0 {
		return pagination.DefaultTraditionalParams()
	}
	return &pagination.TraditionalParams{
		PageNum:  params.PageNum,
		PageSize: params.PageSize,
	}
}
func (params *NftAssetSearch) GetQuickPageParams() *pagination.PageQuickParams {
	if params.PageParams == nil {
		return pagination.DefaultQuickParams()
	}
	return params.PageParams
}
