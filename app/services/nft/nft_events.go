package nft

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/search"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	eventListNum = 50
	eventLastID  primitive.ObjectID
)

type (
	NftEventInput struct {
		*search.NftEventSearch
	}
)

func PageNftEvent(ctx context.Context, currentUID uint64, in *NftEventInput) ([]*models.NftEvent, pagination.Pagination, error) {
	in.NftEventSearch.EventTypes = []string{"successful", "transfer"}
	in.NftEventSearch.SortBy = "created_date_desc"
	ctxWithUID := context.WithValue(ctx, utils.CurrentUIDKey{}, currentUID)
	return models.QuickPageNftEvent(ctxWithUID, in.NftEventSearch)
}

func InitNftEvent(ctx context.Context) error {

	c, err := models.CountNftAsset(ctx, &search.NftAssetSearch{})
	if err != nil {
		fmt.Println("count nft_asset error: ", err.Error())
		return err
	}
	if c == 0 {
		return nil
	}
	times := int(math.Ceil(float64(c) / float64(eventListNum)))
	for i := 0; i < times; i++ {
		err := initNftEvent(ctx)
		if err != nil {
			fmt.Println("init err: ", err.Error())
		}

	}
	return nil
}
func initNftEvent(ctx context.Context) error {
	lists, err := models.NewListNftAsset(ctx, &search.NftAssetSearch{ListNum: int64(eventListNum), LastID: eventLastID})
	if err != nil {
		return err
	}
	for _, v := range lists {
		err := updateNftAssetOneEvent(ctx, v)
		if err != nil {
			fmt.Println("init nft_event one err: ", err.Error())
		}
	}
	eventLastID = lists[len(lists)-1].ID
	return nil
}
func updateNftAssetOneEvent(ctx context.Context, asset *models.NftAsset) error {
	if asset == nil {
		return errors.New("updateNftAssetOneEvent asset is nil")
	}
	if asset.Asset.AssetContract == nil {
		return errors.New("updateNftAssetOneEvent asset_contract is nil")
	}
	params := &OpensaeInput{
		AssetContractAddress: asset.Asset.AssetContract.Address,
		TokenId:              asset.TokenId,
	}
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second * 1)
		out, err := ListEventOut(ctx, params)
		if err != nil {
			fmt.Println("list err: ", err.Error())
			return err
		}
		err = updateNftEvent(ctx, asset, out.AssetEvents)
		if err != nil {
			fmt.Println("update nft_event err: ", err.Error())
			return err
		}
		if out.Next == "" {
			break
		}
		params.Cursor = out.Previous
	}
	return nil
}

func updateNftEvent(ctx context.Context, asset *models.NftAsset, events []*models.AssetEvent) error {
	if asset == nil {
		return nil
	}
	for _, event := range events {
		event.NftAssetID = asset.ID
		err := models.SaveNftEvent(ctx, event)
		if err != nil {
			return err
		}
	}
	return nil
}
