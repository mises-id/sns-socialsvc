package nft

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
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
	nft_asset, err := models.FindNftAssetByID(ctxWithUID, id)
	InitUserNftAssetOne(ctx, nft_asset)
	return nft_asset, err
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
		err := SaveUserNftAssets(ctx, uint64(i))
		if err != nil {
			fmt.Printf("uid[%d],err:%s", i, err.Error())
		}
	}
	return nil
}

func InitUserNftAssets(ctx context.Context, uid uint64) error {
	types := enum.NftTagableTypeOwner
	object_id := strconv.Itoa(int(uid))
	params := &search.NftLogSearch{NftTagableType: types, ObjectID: object_id}
	_, err := models.FindNftLog(ctx, params)
	if err == nil {
		return nil
	}
	if err == mongo.ErrNoDocuments {
		err2 := SaveUserNftAssets(ctx, uid)
		if err2 != nil {
			return err2
		}
		return models.CreateNftLog(ctx, types, object_id)
	}
	return err
}
func InitUserNftAssetOne(ctx context.Context, asset *models.NftAsset) error {
	if asset == nil {
		return nil
	}
	types := enum.NftTagableTypeAsset
	object_id := asset.ID.Hex()
	params := &search.NftLogSearch{NftTagableType: types, ObjectID: object_id}
	_, err := models.FindNftLog(ctx, params)
	if err == nil {
		return nil
	}
	if err == mongo.ErrNoDocuments {
		err2 := UpdateNftAssetOne(ctx, asset)
		if err2 != nil {
			return err2
		}
		return models.CreateNftLog(ctx, types, object_id)
	}
	return err
}

func UpdateNftAssetOne(ctx context.Context, asset *models.NftAsset) error {
	if asset == nil {
		return nil
	}
	uid := asset.UID
	params := &SingleAssetInput{
		AssetContractAddress: asset.AssetContract.Address,
		TokenId:              asset.TokenId,
	}
	new_asset, err := GetSingleAssetOut(ctx, params)
	if err != nil {
		return err
	}
	new_asset.Blockchains = asset.Blockchains
	new_asset.UID = asset.UID
	old_owner_address := asset.Owner.Address
	new_owner_address := new_asset.Owner.Address
	//address change  handlers uid nft_avatar
	if old_owner_address != new_owner_address {
		user, err := models.FindUserEthAddress(ctx, uid)
		if err != nil {
			return err
		}
		if user.NftAvatar != nil && user.NftAvatar.NftAssetID == asset.ID {
			user.NftAvatar = nil
			err := models.UpdateUserAvatar(ctx, user)
			if err != nil {
				return err
			}
		}
		new_owner_user, _ := models.FindUserByEthAddress(ctx, new_owner_address)
		if new_owner_user != nil {
			new_asset.UID = new_owner_user.UID
		}
	}

	if err = models.SaveNftAsset(ctx, new_asset); err != nil {
		return err
	}
	return nil
}

func SaveUserNftAssets(ctx context.Context, uid uint64) error {
	user, err := models.FindUserEthAddress(ctx, uid)
	if err != nil {
		return err
	}
	address := user.EthAddress
	if address == "" {
		return errors.New("invalid address")
	}
	return saveUserNftAssets(ctx, uid, address)
}

func saveUserNftAssets(ctx context.Context, uid uint64, address string) error {
	params := &ListAssetInput{
		Owner: address,
	}
	for i := 0; i < 100; i++ {
		out, err := ListAssetOut(ctx, params)
		if err != nil {
			return err
		}
		err = updateOrCreateNftAssets(ctx, uid, out.Assets)
		if err != nil {
			return err
		}
		if out.Next == "" {
			break
		}
		params.Cursor = out.Previous
	}
	return nil
}

func updateOrCreateNftAssets(ctx context.Context, uid uint64, assets []*models.Asset) error {
	for _, asset := range assets {
		asset.UID = uid
		asset.Blockchains = enum.EthMain
		err := models.SaveNftAsset(ctx, asset)
		if err != nil {
			return err
		}
	}
	return nil
}
