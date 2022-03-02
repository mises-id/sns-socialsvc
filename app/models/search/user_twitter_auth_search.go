package search

import (
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type (
	UserTwitterAuthSearch struct {
		UID            uint64
		UIDs           []uint64
		Misesid        string
		Misesids       []string
		TwitterUserId  string
		TwitterUserIds []string
		//sort
		//limit
		ListNum int64
		//page
		PageNum  int64 `json:"page_num" query:"page_num"`
		PageSize int64 `json:"page_size" query:"page_size"`
	}
)

func (params *UserTwitterAuthSearch) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//base

	//where
	if params.UID != 0 {
		params.UIDs = []uint64{params.UID}
	}
	if params.UIDs != nil && len(params.UIDs) > 0 {
		chain = chain.Where(bson.M{"uid": bson.M{"$in": params.UIDs}})
	}
	if params.Misesid != "" {
		params.Misesids = []string{params.Misesid}
	}
	if params.Misesids != nil && len(params.Misesids) > 0 {
		chain = chain.Where(bson.M{"misesid": bson.M{"$in": params.Misesids}})
	}
	if params.TwitterUserId != "" {
		params.TwitterUserIds = []string{params.TwitterUserId}
	}
	if params.TwitterUserIds != nil && len(params.TwitterUserIds) > 0 {
		chain = chain.Where(bson.M{"twitter_user_id": bson.M{"$in": params.TwitterUserIds}})
	}
	//sort
	chain = chain.Sort(bson.M{"_id": -1})
	//limit
	if (params.PageNum <= 0 || params.PageSize <= 0) && params.ListNum > 0 {
		chain = chain.Limit(params.ListNum)
	}
	return chain
}

func (params *UserTwitterAuthSearch) GetPageParams() *pagination.TraditionalParams {
	if params.PageNum <= 0 || params.PageSize <= 0 {
		return pagination.DefaultTraditionalParams()
	}
	return &pagination.TraditionalParams{
		PageNum:  params.PageNum,
		PageSize: params.PageSize,
	}
}
