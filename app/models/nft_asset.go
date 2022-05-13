package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Asset struct {
		UID               uint64           `bson:"uid"`
		Blockchains       enum.Blockchains `bson:"blockchains"`
		AssetId           int64            `bson:"asset_id" json:"id"`
		TokenId           string           `json:"token_id" bson:"token_id"`
		NumSales          int64            `json:"num_sales" bson:"num_sales"`
		ImageURL          string           `json:"image_url" bson:"image_url"`
		ImagePreviewUrl   string           `json:"image_preview_url" bson:"image_preview_url"`
		ImageThumbnailUrl string           `json:"image_thumbnail_url" bson:"image_thumbnail_url"`
		AnimationUrl      string           `json:"animation_url" bson:"animation_url"`
		BackgroundColor   string           `json:"background_color" bson:"background_color"`
		Name              string           `json:"name" bson:"name"`
		Description       string           `json:"description" bson:"description"`
		ExternalLink      string           `json:"external_link" bson:"external_link"`
		AssetContract     *AssetContract   `json:"asset_contract" bson:"asset_contract"`
		PermaLink         string           `json:"permalink" bson:"perma_link"`
		Collection        *NftCollection   `json:"collection" bson:"collection"`
		Decimals          int64            `json:"decimals" bson:"decimals"`
		TokenMetaData     string           `json:"token_meta_data" bson:"token_meta_data"`
		Owner             Account          `json:"owner" bson:"owner"`
		Creator           Account          `json:"creator" bson:"creator"`
		LastSale          *Sale            `json:"last_sale" bson:"last_sale"`
		ListingDate       string           `json:"listing_date" bson:"listing_date"`
		IsPresale         bool             `json:"is_presale" bson:"is_presale"`
		UpdatedAt         time.Time        `bson:"updated_at,omitempty"`
	}
)
type SaleInfo struct {
	SaleState string `bson:"sale_state"`
	Symbol    string `bson:"symbol"`
	Prices    string `bson:"prices"`
}
type NftAsset struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Asset         `bson:"inline"`
	SaleInfo      *SaleInfo
	CommentsCount uint64    `bson:"comments_count,omitempty"`
	LikesCount    uint64    `bson:"likes_count,omitempty"`
	ForwardsCount uint64    `bson:"forwards_count,omitempty"`
	User          *User     `bson:"-"`
	IsLiked       bool      `bson:"-"`
	CreatedAt     time.Time `bson:"created_at,omitempty"`
}

func SaveNftAsset(ctx context.Context, data *Asset) error {
	data.UpdatedAt = time.Now()
	opt := &options.FindOneAndUpdateOptions{}
	opt.SetUpsert(true)
	opt.SetReturnDocument(1)
	result := db.DB().Collection("nftassets").FindOneAndUpdate(ctx, bson.M{"asset_contract.address": data.AssetContract.Address, "token_id": data.TokenId}, bson.D{{Key: "$set", Value: data}}, opt)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func CountNftAsset(ctx context.Context, params IAdminParams) (int64, error) {

	var res int64
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Model(&NftAsset{}).Count(&res).Error
	if err != nil {
		return res, err
	}

	return res, nil
}

func FindNftAsset(ctx context.Context, params IAdminParams) (*NftAsset, error) {

	res := &NftAsset{}
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Get(res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}
func FindNftAssetByID(ctx context.Context, id primitive.ObjectID) (*NftAsset, error) {

	res := &NftAsset{}
	err := db.ODM(ctx).First(res, bson.M{"_id": id}).Error
	if err != nil {
		return nil, err
	}
	return res, PreloadNftAssets(ctx, res)
}

func FindNftAssetByIDs(ctx context.Context, ids ...primitive.ObjectID) ([]*NftAsset, error) {
	res := make([]*NftAsset, 0)
	err := db.ODM(ctx).Where(bson.M{"_id": bson.M{"$in": ids}}).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, PreloadNftAssets(ctx, res...)
}

func PreloadNftAssets(ctx context.Context, assets ...*NftAsset) error {
	var err error
	if err = preloadAssetsUser(ctx, assets...); err != nil {
		return err
	}
	if err = preloadAssetsLikeState(ctx, assets...); err != nil {
		return err
	}
	return nil
}

func preloadAssetsUser(ctx context.Context, assets ...*NftAsset) error {
	userIds := make([]uint64, 0)
	for _, v := range assets {
		userIds = append(userIds, v.UID)
	}
	users, err := FindUserByIDs(ctx, userIds...)
	if err != nil {
		return err
	}
	userMap := make(map[uint64]*User)
	for _, user := range users {
		userMap[user.UID] = user
	}
	for _, v := range assets {
		v.User = userMap[v.UID]
	}
	return nil
}

func preloadAssetsLikeState(ctx context.Context, assets ...*NftAsset) error {
	currentUID, ok := ctx.Value(utils.CurrentUIDKey{}).(uint64)
	if !ok || currentUID == 0 {
		return nil
	}
	ids := make([]primitive.ObjectID, len(assets))
	for i, v := range assets {
		ids[i] = v.ID
	}
	likeMap, err := GetLikeMap(ctx, currentUID, ids, enum.LikeNft, false)
	if err != nil {
		return err
	}
	for _, v := range assets {
		v.IsLiked = likeMap[v.ID] != nil
	}
	return nil
}

func NewListNftAsset(ctx context.Context, params IAdminParams) ([]*NftAsset, error) {
	res := make([]*NftAsset, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	err := chain.Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func QuickPageNftAsset(ctx context.Context, params IAdminQuickPageParams) ([]*NftAsset, pagination.Pagination, error) {
	out := make([]*NftAsset, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	pageParams := params.GetQuickPageParams()
	paginator := pagination.NewQuickPaginator(pageParams.Limit, pageParams.NextID, chain)
	page, err := paginator.Paginate(&out)
	if err != nil {
		return nil, nil, err
	}

	return out, page, PreloadNftAssets(ctx, out...)
}

func (s *NftAsset) IncNftAssetCounter(ctx context.Context, counterKey string, values ...int) error {
	if counterKey == "" {
		return nil
	}
	value := 1
	if len(values) > 0 {
		value = values[0]
	}
	if value < 0 {
		switch counterKey {
		case "likes_count":
			if int(s.LikesCount)+value < 0 {
				value = -int(s.LikesCount)
			}
		case "comments_count":
			if int(s.CommentsCount)+value < 0 {
				value = -int(s.CommentsCount)
			}
		}
	}
	return db.DB().Collection("nftassets").FindOneAndUpdate(ctx, bson.M{"_id": s.ID},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   counterKey,
				Value: value,
			}}},
		}).Err()
}
