package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NftCollection struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty"`
	PrimaryAssetContracts   []*AssetContract   `json:"primary_asset_contracts" bson:"primary_asset_contracts"`
	Stats                   *Stats             `json:"stats" bson:"stats"`
	PaymentToken            []*PaymentToken    `json:"payment_token" bson:"payment_token"           `
	BannerImageUrl          string             `json:"banner_image_url" bson:"banner_image_url"`
	ChatUrl                 string             `json:"chat_url" bson:"chat_url"`
	CreatedDate             string             `json:"created_date" bson:"created_date"`
	DefaultToFiat           bool               `json:"default_to_fiat" bson:"default_to_fiat"`
	Description             string             `json:"description" bson:"description"`
	DevBuyerFeeBasisPoints  string             `json:"dev_buyer_fee_basis_points" bson:"dev_buyer_fee_basis_points"`
	DevSellerFeeBasisPoints string             `json:"dev_seller_fee_basis_points" bson:"dev_seller_fee_basis_points"`
	DiscordUrl              string             `json:"discord_url" bson:"discord_url"`
	//DisplayData                 DisplayData        `json:"display_data" bson:"display_data"`
	ExternalUrl                 string `json:"external_url" bson:"external_url"`
	Featured                    bool   `json:"featured" bson:"featured"`
	FeaturedImageUrl            string `json:"featured_image_url" bson:"featured_image_url"`
	Hidden                      bool   `json:"hidden" bson:"hidden"`
	SafelistRequestStatus       string `json:"safelist_request_status" bson:"safelist_request_status"`
	ImageURL                    string `json:"image_url" bson:"image_url"`
	IsSubjectToWhitelist        bool   `json:"is_subject_to_whitelist" bson:"is_subject_to_whitelist"`
	LargeImageUrl               string `json:"large_image_url" bson:"large_image_url"`
	MediumUsername              string `json:"medium_username" bson:"medium_username"`
	Name                        string `json:"name" bson:"name"`
	OnlyProxiedTransfers        bool   `json:"only_proxied_transfers" bson:"only_proxied_transfers"`
	OpenseaBuyerFeeBasisPoints  string `json:"opensea_buyer_fee_basis_points" bson:"opensea_buyer_fee_basis_points"`
	OpenseaSellerFeeBasisPoints string `json:"opensea_seller_fee_basis_points" bson:"opensea_seller_fee_basis_points"`
	PayoutAddress               string `json:"payout_address" bson:"payout_address"`
	RequireEmail                bool   `json:"require_email" bson:"require_email"`
	ShortDescription            string `json:"short_description" bson:"short_description"`
	Slug                        string `json:"slug" bson:"slug"`
	TelegramUrl                 string `json:"telegram_url" bson:"telegram_url"`
	TwitterUsername             string `json:"twitter_username" bson:"twitter_username"`
	InstagramUsername           string `json:"instagram_username" bson:"instagram_username"`
	WikiUrl                     string `json:"wiki_url" bson:"wiki_url"`
	OwnedAssetsCount            int64  `json:"owned_asset_count" bson:"owned_assets_count"`
	//PaymentTokens               []PaymentToken     `json:"payment_tokens" bson:"payment_tokens"`
	Traits interface{} `json:"traits" bson:"traits"`
}
