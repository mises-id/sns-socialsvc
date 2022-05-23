package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AssetEvent struct {
	ApprovedAccount    *Account           `json:"approved_account" bson:"approved_account"`
	AssetBundle        string             `json:"asset_bundle" bson:"asset_bundle"`
	AuctionType        string             `json:"auction_type" bson:"auction_type"`
	BidAmount          string             `json:"bid_amount" bson:"bid_amount"`
	CollectionSlug     string             `json:"collection_slug" bson:"collection_slug"`
	ContractAddress    string             `json:"contract_address" bson:"contract_address"`
	CreatedDate        string             `json:"created_date" bson:"created_date"`
	CustomEventName    string             `json:"custom_event_name" bson:"custom_event_name"`
	DevFeePaymentEvent string             `json:"dev_fee_payment_event" bson:"dev_fee_payment_event"`
	Duration           string             `json:"duration" bson:"duration"`
	EndingPrice        string             `json:"ending_price" bson:"ending_price"`
	EventType          string             `json:"event_type" bson:"event_type"`
	FromAccount        *Account           `json:"from_account" bson:"from_account"`
	AssetEventId       int64              `json:"id" bson:"asset_event_id"`
	IsPrivate          bool               `json:"is_private" bson:"is_private"`
	OwnerAccount       *Account           `json:"owner_account" bson:"owner_account"`
	PaymentToken       *PaymentToken      `json:"payment_token" bson:"payment_token"`
	Quantity           string             `json:"quantity" bson:"quantity"`
	Seller             *Account           `json:"seller" bson:"seller"`
	StartingPrice      string             `json:"starting_price" bson:"starting_price"`
	ToAccount          *Account           `json:"to_account" bson:"to_account"`
	TotalPrice         string             `json:"total_price" bson:"total_price"`
	Transaction        *Transaction       `json:"transaction" bson:"transaction"`
	WinnerAccount      *Account           `json:"winner_account" bson:"winner_account"`
	NftAssetID         primitive.ObjectID `bson:"nft_asset_id"`
	UpdatedAt          time.Time          `bson:"updated_at,omitempty"`
}
type NftEvent struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	AssetEvent `bson:"inline"`
}

func SaveNftEvent(ctx context.Context, data *AssetEvent) error {
	data.UpdatedAt = time.Now()
	opt := &options.FindOneAndUpdateOptions{}
	opt.SetUpsert(true)
	opt.SetReturnDocument(1)
	result := db.DB().Collection("nftevents").FindOneAndUpdate(ctx, bson.M{"asset_event_id": data.AssetEventId, "nft_asset_id": data.NftAssetID}, bson.D{{Key: "$set", Value: data}}, opt)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func QuickPageNftEvent(ctx context.Context, params IAdminQuickPageParams) ([]*NftEvent, pagination.Pagination, error) {
	out := make([]*NftEvent, 0)
	chain := params.BuildAdminSearch(db.ODM(ctx))
	pageParams := params.GetQuickPageParams()
	paginator := pagination.NewQuickPaginator(pageParams.Limit, pageParams.NextID, chain)
	page, err := paginator.Paginate(&out)
	if err != nil {
		return nil, nil, err
	}

	return out, page, PreloadNftEvents(ctx, out...)
}
func PreloadNftEvents(ctx context.Context, events ...*NftEvent) error {
	if err := preloadNftAccountMisesUser(ctx, events...); err != nil {
		return err
	}
	return nil
}

func preloadNftAccountMisesUser(ctx context.Context, events ...*NftEvent) error {
	eth_addresses := make([]string, 0)
	for _, event := range events {
		if event.FromAccount != nil && event.FromAccount.Address != "" {
			eth_addresses = append(eth_addresses, event.FromAccount.Address)
		}
		if event.ToAccount != nil && event.ToAccount.Address != "" {
			eth_addresses = append(eth_addresses, event.ToAccount.Address)
		}
	}
	addressMap, err := GetUserMapByEthAddresses(ctx, eth_addresses...)
	if err != nil {
		return err
	}
	for _, event := range events {
		if event.FromAccount != nil && event.FromAccount.Address != "" {
			event.FromAccount.MisesUser = addressMap[event.FromAccount.Address]
		}
		if event.ToAccount != nil && event.ToAccount.Address != "" {
			event.ToAccount.MisesUser = addressMap[event.ToAccount.Address]
		}
	}
	return nil
}
