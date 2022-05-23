package nft

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"github.com/sirupsen/logrus"
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
		err := UpdateUserNftAssets(ctx, uint64(i))
		if err != nil {
			fmt.Printf("uid[%d],err:%s", i, err.Error())
		}
	}
	return nil
}

func SaveUserNftAsset(ctx context.Context, uid uint64, assets []*models.Asset) error {
	if err := SaveUserNftLog(ctx, uid); err != nil {
		return err
	}
	return updateOrCreateNftAssets(ctx, uid, assets)
}

func SaveUserNftLog(ctx context.Context, uid uint64) error {
	types := enum.NftTagableTypeOwner
	object_id := strconv.Itoa(int(uid))

	/* params := &search.NftLogSearch{NftTagableType: types, ObjectID: object_id}
	_, err := models.FindNftLog(ctx, params)
	if err == nil {
		return nil
	}
	if err == mongo.ErrNoDocuments {
		err2 := UpdateUserNftAssets(ctx, uid)
		if err2 != nil {
			return err2
		}
		return models.CreateNftLog(ctx, types, object_id)
	} */
	return models.SaveNftLog(ctx, types, object_id)
}

func UpdateOpenseaNft(ctx context.Context) error {
	sh, _ := time.ParseDuration("-8h")
	st := time.Now().UTC().Add(sh)
	logs, err := models.ListNftLog(ctx, &search.NftLogSearch{
		NeedUpdate:     true,
		UpdatedAt:      &st,
		NftTagableType: enum.NftTagableTypeOwner,
		ListNum:        10,
	})
	if err != nil {
		return err
	}
	for _, log := range logs {
		uid, err := strconv.ParseUint(log.ObjectID, 10, 64)
		if err != nil {
			continue
		}
		err = UpdateUserNftAssets(ctx, uid)
		if err != nil {
			fmt.Println("update user nft assets error: ", err.Error())
			continue
		}
		log.ForceUpdate = false
		models.UpdateNftLog(ctx, log)
	}
	return nil
}

//update nft assets by uid
func UpdateUserNftAssets(ctx context.Context, uid uint64) error {
	user, err := models.FindUserEthAddress(ctx, uid)
	if err != nil {
		return err
	}
	address := user.EthAddress
	if address == "" {
		return errors.New("invalid address")
	}
	return updateUserNftAssets(ctx, uid, address)
}

func updateUserNftAssets(ctx context.Context, uid uint64, address string) error {
	params := &ListAssetInput{
		Owner: address,
	}
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second * 1)
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
		nft_asset, err := saveUserNftAssetOne(ctx, uid, asset)
		if err != nil {
			logrus.Printf("save user nft_asset one err: %s", err.Error())
			return err
		}
		//save event
		if err = updateNftAssetOneEvent(ctx, nft_asset); err != nil {
			logrus.Printf("update nft_asset one event err: %s", err.Error())
			return err
		}
	}
	return nil
}

func saveUserNftAssetOne(ctx context.Context, uid uint64, asset *models.Asset) (*models.NftAsset, error) {
	time.Sleep(time.Second * 1)
	if asset == nil {
		return nil, errors.New("saveUserNftAssetOne asset is nil")
	}
	params := &SingleAssetInput{
		AssetContractAddress: asset.AssetContract.Address,
		TokenId:              asset.TokenId,
	}
	new_asset, err := GetSingleAssetOut(ctx, params)
	if err != nil {
		return nil, err
	}
	nft_asset, err := models.FindNftAssetByAddressAndToken(ctx, asset.AssetContract.Address, asset.TokenId)
	//create
	if err != nil && err == mongo.ErrNoDocuments {
		create := &models.NftAsset{Asset: *new_asset}
		create.UID = uid
		create.Blockchains = enum.EthMain
		nft_asset, err = models.CreateNftAsset(ctx, create)
		return nft_asset, err
	}
	//update
	err = updateNftAssetOne(ctx, nft_asset, new_asset)
	return nft_asset, err
}

func updateNftAssetOne(ctx context.Context, nft_asset *models.NftAsset, new_asset *models.Asset) error {
	if nft_asset == nil {
		return errors.New("updateNftAssetOne nft_asset is nil")
	}
	if new_asset.AssetContract == nil {
		return errors.New("updateNftAssetOne new_asset is nil")
	}
	uid := nft_asset.UID
	new_asset.Blockchains = nft_asset.Blockchains
	new_asset.UID = nft_asset.UID
	old_owner_address := nft_asset.Owner.Address
	new_owner_address := new_asset.Owner.Address
	//address change  handlers uid nft_avatar
	if old_owner_address != new_owner_address {
		user, err := models.FindUserEthAddress(ctx, uid)
		if err != nil {
			return err
		}
		if user.NftAvatar != nil && user.NftAvatar.NftAssetID == nft_asset.ID {
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
	nft_asset.Asset = *new_asset
	if err := models.UpdateNftAsset(ctx, nft_asset); err != nil {
		return err
	}

	return nil
}
