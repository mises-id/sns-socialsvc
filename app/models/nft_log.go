package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	NftLog struct {
		ID             primitive.ObjectID  `bson:"_id,omitempty"`
		NftTagableType enum.NftTagableType `bson:"nft_tagable_type"`
		ObjectID       string              `bson:"object_id"`
		ForceUpdate    bool                `bson:"force_update"`
		Num            uint64              `bson:"num"`
		UpdatedAt      time.Time           `bson:"updated_at,omitempty"`
		CreatedAt      time.Time           `bson:"created_at,omitempty"`
		UpdateType     string              `bson:"-"`
	}
)

func CreateNftLog(ctx context.Context, nft_tagable_type enum.NftTagableType, object_id string) error {
	created := bson.M{}
	created["created_at"] = time.Now()
	created["updated_at"] = time.Now()
	created["nft_tagable_type"] = nft_tagable_type
	created["object_id"] = object_id
	opt := &options.FindOneAndUpdateOptions{}
	opt.SetUpsert(true)
	opt.SetReturnDocument(1)
	result := db.DB().Collection("nftlogs").FindOneAndUpdate(ctx,
		bson.M{"nft_tagable_type": nft_tagable_type, "object_id": object_id},
		bson.D{{Key: "$set", Value: created}}, opt)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func UpdateNftLog(ctx context.Context, log *NftLog) error {

	update := bson.M{}
	update["updated_at"] = time.Now()
	update["force_update"] = log.ForceUpdate
	if log.UpdateType == "update" {
		update["num"] = log.Num
	}
	_, err := db.DB().Collection("nftlogs").UpdateByID(ctx, log.ID, bson.D{{Key: "$set", Value: update}})
	return err
}

func SaveNftLog(ctx context.Context, nft_tagable_type enum.NftTagableType, object_id string) error {
	types := enum.NftTagableTypeOwner
	params := &search.NftLogSearch{NftTagableType: types, ObjectID: object_id}
	log, err := FindNftLog(ctx, params)
	if err == nil {
		log.ForceUpdate = true
		return UpdateNftLog(ctx, log)
	}
	if err == mongo.ErrNoDocuments {
		//create
		return CreateNftLog(ctx, nft_tagable_type, object_id)
	}
	return err
}

func FindNftLog(ctx context.Context, params IAdminParams) (*NftLog, error) {

	res := &NftLog{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Get(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ListNftLog(ctx context.Context, params IAdminParams) ([]*NftLog, error) {

	res := make([]*NftLog, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

//page user
func AdminPageNftLog(ctx context.Context, params IAdminPageParams) ([]*NftLog, pagination.Pagination, error) {
	users := make([]*NftLog, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	pageParams := params.GetPageParams()
	paginator := pagination.NewTraditionalPaginatorAdmin(pageParams.PageNum, pageParams.PageSize, chain)
	page, err := paginator.Paginate(&users)
	if err != nil {
		return nil, nil, err
	}

	return users, page, nil
}
