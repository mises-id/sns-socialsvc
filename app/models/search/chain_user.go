package search

import (
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	ChainUserSearch struct {
		Misesid  string
		Misesids []string
		Status   enum.ChainUserStatus
		Statuses []enum.ChainUserStatus
		TxID     string
		NotTxID  bool
		//sort
		SortKey  string
		SortType enum.SortType
		//limit
		ListNum int64
		//page
		PageNum  int64 `json:"page_num" query:"page_num"`
		PageSize int64 `json:"page_size" query:"page_size"`
	}
)

func (params *ChainUserSearch) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//base
	//where

	if params.Misesid != "" {
		params.Misesids = []string{params.Misesid}
	}
	if params.Misesids != nil && len(params.Misesids) > 0 {
		chain = chain.Where(bson.M{"misesid": bson.M{"$in": params.Misesids}})
	}
	if params.Statuses != nil && len(params.Statuses) > 0 {
		chain = chain.Where(bson.M{"status": bson.M{"$in": params.Statuses}})
	}
	if params.Status > -1 {
		chain = chain.Where(bson.M{"status": params.Status})
	}
	if params.TxID != "" {
		chain = chain.Where(bson.M{"tx_id": params.TxID})
	}
	if params.NotTxID {
		chain = chain.Where(bson.M{"tx_id": ""})
	}
	//sort
	if params.SortKey != "" && params.SortType != 0 {
		chain = chain.Sort(bson.M{params.SortKey: params.SortType})
	}
	chain = chain.Sort(bson.M{"_id": -1})
	//limit
	if (params.PageNum <= 0 || params.PageSize <= 0) && params.ListNum > 0 {
		chain = chain.Limit(params.ListNum)
	}
	return chain
}

func (params *ChainUserSearch) GetPageParams() *pagination.TraditionalParams {
	if params.PageNum <= 0 || params.PageSize <= 0 {
		return pagination.DefaultTraditionalParams()
	}
	return &pagination.TraditionalParams{
		PageNum:  params.PageNum,
		PageSize: params.PageSize,
	}
}
