package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NftAsset struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty"`
	UID                     uint64             `bson:"uid"`
	AssetId                 int64              `bson:"asset_id" json:"id"`
	TokenId                 string             `json:"token_id" bson:"token_id"`
	NumSales                int64              `json:"num_sales" bson:"num_sales"`
	ImageURL                string             `json:"image_url" bson:"image_url"`
	ImagePreviewUrl         string             `json:"image_preview_url" bson:"image_preview_url"`
	ImageThumbnailUrl       string             `json:"image_thumbnail_url" bson:"image_thumbnail_url"`
	AnimationUrl            string             `json:"animation_url" bson:"animation_url"`
	BackgroundColor         string             `json:"background_color" bson:"background_color"`
	Name                    string             `json:"name" bson:"name"`
	Description             string             `json:"description" bson:"description"`
	ExternalLink            string             `json:"external_link" bson:"external_link"`
	AssetContract           AssetContract      `json:"asset_contract" bson:"asset_contract"`
	PermaLink               string             `json:"permalink" bson:"perma_link"`
	Collection              NftCollection      `json:"collection" bson:"collection"`
	Decimals                float64            `json:"decimals" bson:"decimals"`
	TokenMetaData           string             `json:"token_meta_data" bson:"token_meta_data"`
	Owner                   Account            `json:"owner" bson:"owner"`
	Creator                 Account            `json:"user" bson:"creator"`
	Traits                  []Trait            `json:"traits" bson:"traits"`
	LastSale                Sale               `json:"last_sale" bson:"last_sale"`
	ListingDate             string             `json:"listing_date" bson:"listing_date"`
	IsPresale               bool               `json:"is_presale" bson:"is_presale"`
	TransferFeePaymentToken PaymentToken       `json:"transfer_fee_payment_token" bson:"transfer_fee_payment_token"`
	TransferFee             string             `json:"transfer_fee" bson:"transfer_fee"`
	TopBid                  string             `json:"top_bid" bson:"top_bid"`
	SupportsWyvern          bool               `json:"supports_wyvern" bson:"supports_wyvern"`
	TopOwnerships           []Ownership        `json:"top_ownerships" bson:"top_ownerships"`
	Ownership               Ownership          `json:"ownership" bson:"ownership"`
}
