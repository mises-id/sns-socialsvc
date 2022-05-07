package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NftAssetEvent struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	ApprovedAccount    Account            `json:"approved_account" bson:"approved_account"`
	Asset              NftAsset           `json:"asset" bson:"asset"`
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
	FromAccount        Account            `json:"from_account" bson:"from_account"`
	AssetEventId       int64              `json:"id" bson:"asset_event_id"`
	IsPrivate          bool               `json:"is_private" bson:"is_private"`
	OwnerAccount       Account            `json:"owner_account" bson:"owner_account"`
	PaymentToken       PaymentToken       `json:"payment_token" bson:"payment_token"`
	Quantity           string             `json:"quantity" bson:"quantity"`
	Seller             Account            `json:"seller" bson:"seller"`
	StartingPrice      string             `json:"starting_price" bson:"starting_price"`
	ToAccount          Account            `json:"to_account" bson:"to_account"`
	TotalPrice         string             `json:"total_price" bson:"total_price"`
	Transaction        Transaction        `json:"transaction" bson:"transaction"`
	WinnerAccount      Account            `json:"winner_account" bson:"winner_account"`
}
