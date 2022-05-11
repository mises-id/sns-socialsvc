package nft

import (
	"context"
	"errors"
	"fmt"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/app/services/opensea_api"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	NftAssetInput struct {
		*search.NftAssetSearch
	}
	PageLikeParams struct {
		*models.LikeSearch
	}
)

func LikeNftAsset(ctx context.Context, uid uint64, nft_asset_id primitive.ObjectID) (*models.Like, error) {
	ctx = context.WithValue(ctx, utils.CurrentUIDKey{}, uid)
	nft_asset, err := models.FindNftAssetByID(ctx, nft_asset_id)
	if err != nil {
		return nil, err
	}
	like, err := models.FindLike(ctx, uid, nft_asset_id, enum.LikeNft)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == nil {
		return like, nil
	}
	like, err = models.CreateLike(ctx, uid, nft_asset.UID, nft_asset_id, enum.LikeNft)
	if err != nil {
		return nil, err
	}
	return like, nft_asset.IncNftAssetCounter(ctx, "likes_count")
}

func UnlikeNftAsset(ctx context.Context, uid uint64, nft_asset_id primitive.ObjectID) error {
	like, err := models.FindLike(ctx, uid, nft_asset_id, enum.LikeNft)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	nft_asset, err := models.FindNftAssetByID(ctx, nft_asset_id)
	if err != nil {
		return err
	}
	if err = models.DeleteLike(ctx, like.ID); err != nil {
		return err
	}
	return nft_asset.IncNftAssetCounter(ctx, "likes_count", -1)
}

func FindNftAsset(ctx context.Context, currentUID uint64, id primitive.ObjectID) (*models.NftAsset, error) {
	ctxWithUID := context.WithValue(ctx, utils.CurrentUIDKey{}, currentUID)
	return models.FindNftAssetByID(ctxWithUID, id)
}

func PageNftAsset(ctx context.Context, currentUID uint64, in *NftAssetInput) ([]*models.NftAsset, pagination.Pagination, error) {
	ctxWithUID := context.WithValue(ctx, utils.CurrentUIDKey{}, currentUID)
	params := in.NftAssetSearch
	return models.QuickPageNftAsset(ctxWithUID, params)
}

func PageLike(ctx context.Context, currentUID uint64, params *PageLikeParams) ([]*models.Like, pagination.Pagination, error) {
	ctxWithUID := context.WithValue(ctx, utils.CurrentUIDKey{}, currentUID)
	return models.PageLike(ctxWithUID, enum.LikeNft, params.LikeSearch)
}

func Run(ctx context.Context) error {
	for i := 1; i < 100; i++ {
		err := InitUserNftAssets(ctx, uint64(i))
		if err != nil {
			fmt.Printf("uid[%d],err:%s", i, err.Error())
		}
	}
	return nil
}

func InitUserNftAssets(ctx context.Context, uid uint64) error {
	err := models.UpdateUserExtNftState(ctx, uid, true)
	if err != nil {
		return err
	}
	user_ext, err := models.FindUserExt(ctx, uid)
	if err != nil {
		return err
	}
	if !user_ext.NftState {
		return errors.New("nft state false")
	}
	address := user_ext.EthAddress
	if address == "" {
		return errors.New("invalid address")
	}
	return initUserNftAssets(ctx, uid, address)
}

func initUserNftAssets(ctx context.Context, uid uint64, address string) error {
	params := &opensea_api.ListAssetInput{
		Owner: address,
	}
	for i := 0; i < 100; i++ {
		out, err := opensea_api.ListAssetOut(ctx, params)
		if err != nil {
			fmt.Println("list err: ", err.Error())
			return err
		}
		err = updateNftAssets(ctx, uid, out.Assets)
		if err != nil {
			fmt.Println("update err: ", err.Error())
			return err
		}
		if out.Next == "" {
			fmt.Println("no asset")
			break
		}
		params.Cursor = out.Previous
	}
	return nil
}

func updateNftAssets(ctx context.Context, uid uint64, assets []*models.Asset) error {
	for _, asset := range assets {
		asset.UID = uid
		asset.Blockchains = "eth main"
		err := models.SaveNftAsset(ctx, asset)
		if err != nil {
			return err
		}
	}
	return nil
}
