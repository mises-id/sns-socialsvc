package search

import (
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db/odm"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ChannelUserSearch struct {
		ID              primitive.ObjectID
		ChannelID       primitive.ObjectID
		ChannelMisesid  string
		ChannelMisesids []string
		TxID            string
		UID             uint64
		UIDs            []uint64
		//ValidState      enum.UserValidState
		ValidStates []enum.UserValidState
		//AirdropState    enum.ChannelAirdropState
		AirdropStates []enum.ChannelAirdropState
		StartTime     *time.Time `json:"start_time" query:"start_time"`
		EndTime       *time.Time `json:"end_time" query:"end_time"`
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

func (params *ChannelUserSearch) BuildAdminSearch(chain *odm.DB) *odm.DB {
	//base
	//where
	if params.ID != primitive.NilObjectID {
		chain = chain.Where(bson.M{"_id": params.ID})
	}
	if params.ChannelID != primitive.NilObjectID {
		chain = chain.Where(bson.M{"channel_id": params.ChannelID})
	}
	if params.ChannelMisesid != "" {
		params.ChannelMisesids = []string{utils.AddMisesidProfix(params.ChannelMisesid)}
	}
	if params.ChannelMisesids != nil && len(params.ChannelMisesids) > 0 {
		chain = chain.Where(bson.M{"channel_misesid": bson.M{"$in": params.ChannelMisesids}})
	}
	if params.UID > 0 {
		params.UIDs = []uint64{params.UID}
	}
	if params.UIDs != nil && len(params.UIDs) > 0 {
		chain = chain.Where(bson.M{"uid": bson.M{"$in": params.UIDs}})
	}
	if params.ValidStates != nil && len(params.ValidStates) > 0 {
		chain = chain.Where(bson.M{"valid_state": bson.M{"$in": params.ValidStates}})
	}
	if params.AirdropStates != nil && len(params.AirdropStates) > 0 {
		chain = chain.Where(bson.M{"airdrop_state": bson.M{"$in": params.AirdropStates}})
	}
	if params.TxID != "" {
		chain = chain.Where(bson.M{"tx_id": params.TxID})
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
	chain = chain.Sort(bson.M{"_id": -1})
	//limit
	if (params.PageNum <= 0 || params.PageSize <= 0) && params.ListNum > 0 {
		chain = chain.Limit(params.ListNum)
	}
	return chain
}

func (params *ChannelUserSearch) GetPageParams() *pagination.TraditionalParams {
	if params.PageNum <= 0 || params.PageSize <= 0 {
		return pagination.DefaultTraditionalParams()
	}
	return &pagination.TraditionalParams{
		PageNum:  params.PageNum,
		PageSize: params.PageSize,
	}
}
