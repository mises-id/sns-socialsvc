package models

import (
	"encoding/json"
)

type Sale struct {
	AssetBundle    string        `json:"asset_bundle" bson:"asset_bundle"`
	EventType      string        `json:"event_type" bson:"event_type"`
	EventTimestamp string        `json:"event_timestamp" bson:"event_timestamp"`
	AuctionType    string        `json:"auction_type" bson:"auction_type"`
	TotalPrice     string        `json:"total_price" bson:"total_price"`
	PaymentToken   *PaymentToken `json:"payment_token" bson:"payment_token"`
	Transaction    *Transaction  `json:"transaction" bson:"transaction"`
	CreatedDate    string        `json:"created_date" bson:"created_date"`
	Quantity       interface{}   `json:"quantity" bson:"quantity"`
}

type PaymentToken struct {
	ID       int64  `json:"id" bson:"id"`
	Symbol   string `json:"symbol,omitempty" bson:"symbol"`
	Address  string `json:"address,omitempty" bson:"address"`
	ImageURL string `json:"image_url,omitempty" bson:"image_url"`
	Name     string `json:"name,omitempty" bson:"name"`
	Decimals int64  `json:"decimals" bson:"decimals"`
	ETHPrice string `json:"eth_price,omitempty" bson:"eth_price"`
	USDPrice string `json:"usd_price,omitempty" bson:"usd_price"`
}

type Trait struct {
	TraitType   string      `json:"trait_type" bson:"trait_type"`
	Value       interface{} `json:"value" bson:"value"`
	DisplayType string      `json:"display_type" bson:"display_type"`
	MaxValue    json.Number `json:"max_value" bson:"max_value"`
	TraitCount  json.Number `json:"trait_count" bson:"trait_count"`
	Order       interface{} `json:"order" bson:"order"`
}

type Ownership struct {
	Owner    *Account `json:"owner" bson:"owner"`
	Quantity string   `json:"quantity" bson:"quantity"`
}

type Account struct {
	Address       string      `json:"address" bson:"address"`
	ProfileImgUrl string      `json:"profile_img_url" bson:"profile_img_url"`
	User          interface{} `json:"user" bson:"user"`
	Config        string      `json:"config" bson:"config"`
	MisesUser     *User       `bson:"-"`
}

type Bundle struct {
	Maker         *Account       `json:"maker" bson:"maker"`
	Slug          string         `json:"slug" bson:"slug"`
	Assets        []*NftAsset    `json:"assets" bson:"assets"`
	Schemas       []string       `json:"schemas" bson:"schemas"`
	Name          string         `json:"name" bson:"name"`
	Description   string         `json:"description" bson:"description"`
	ExternalLink  string         `json:"external_link" bson:"external_link"`
	AssetContract *AssetContract `json:"asset_contract" bson:"asset_contract"`
	PermaLink     string         `json:"perma_link" bson:"perma_link"`
	SellOrders    []*NftOrder    `json:"sell_orders" bson:"sell_orders"`
}

type Stats struct {
	OneDayVolume          float64 `json:"one_day_volume" bson:"one_day_volume"`
	OneDayChange          float64 `json:"one_day_change" bson:"one_day_change"`
	OneDaySales           float64 `json:"one_day_sales" bson:"one_day_sales"`
	OneDayAveragePrice    float64 `json:"one_day_average_price" bson:"one_day_average_price"`
	SevenDayVolume        float64 `json:"seven_day_volume" bson:"seven_day_volume"`
	SevenDayChange        float64 `json:"seven_day_change" bson:"seven_day_change"`
	SevenDaySales         float64 `json:"seven_day_sales" bson:"seven_day_sales"`
	SevenDayAveragePrice  float64 `json:"seven_day_average_price" bson:"seven_day_average_price"`
	ThirtyDayVolume       float64 `json:"thirty_day_volume" bson:"thirty_day_volume"`
	ThirtyDayChange       float64 `json:"thirty_day_change" bson:"thirty_day_change"`
	ThirtyDaySales        float64 `json:"thirty_day_sales" bson:"thirty_day_sales"`
	ThirtyDayAveragePrice float64 `json:"thirty_day_average_price" bson:"thirty_day_average_price"`
	TotalVolume           float64 `json:"total_volume" bson:"total_volume"`
	TotalSales            float64 `json:"total_sales" bson:"total_sales"`
	TotalSupply           float64 `json:"total_supply" bson:"total_supply"`
	Count                 float64 `json:"count" bson:"count"`
	NumOwners             int64   `json:"num_owners" bson:"num_owners"`
	AveragePrice          float64 `json:"average_price" bson:"average_price"`
	NumReports            int64   `json:"num_reports" bson:"num_reports"`
	MarketCap             float64 `json:"market_cap" bson:"market_cap"`
	FloorPrice            string  `json:"floor_price" bson:"floor_price"`
}

type DisplayData struct {
	CardDisplayStyle string `json:"card_display_style" bson:"card_display_style"`
}

type AssetContract struct {
	Address                     string `json:"address" bson:"address"`
	AddressContractType         string `json:"asset_contract_type" bson:"address_contract_type"`
	CreatedDate                 string `json:"created_date" bson:"created_date"`
	Name                        string `json:"name" bson:"name"`
	NftVersion                  string `json:"nft_versiom" bson:"nft_version"`
	OpenseaVersion              string `json:"opensea_version" bson:"opensea_version"`
	Owner                       int    `json:"owner" bson:"owner"`
	SchemaName                  string `json:"schema_name" bson:"schema_name"`
	Symbol                      string `json:"symbol" bson:"symbol"`
	TotalSupply                 string `json:"total_supply" bson:"total_supply"`
	ImageURL                    string `json:"image_url" bson:"image_url"`
	Description                 string `json:"description" bson:"description"`
	ExternalLink                string `json:"external_link" bson:"external_link"`
	DefaultToFiat               bool   `json:"default_to_fiat" bson:"default_to_fiat"`
	DevBuyerFeeBasisPoints      int    `json:"dev_buyer_fee_basis_points" bson:"dev_buyer_fee_basis_points"`
	DevSellerFeeBasisPoints     int    `json:"dev_seller_fee_basis_points" bson:"dev_seller_fee_basis_points"`
	OnlyProxiedTransfers        bool   `json:"only_proxied_transfers" bson:"only_proxied_transfers"`
	OpenseaBuyerFeeBasisPoints  int    `json:"opensea_buyer_fee_basis_points" bson:"opensea_buyer_fee_basis_points"`
	OpenseaSellerFeeBasisPoints int    `json:"opensea_seller_fee_basis_points" bson:"opensea_seller_fee_basis_points"`
	BuyerFeeBasisPoints         int    `json:"buyer_fee_basis_points" bson:"buyer_fee_basis_points"`
	SellerFeeBasisPoints        int    `json:"seller_fee_basis_points" bson:"seller_fee_basis_points"`
	PayoutAddress               string `json:"payout_address" bson:"payout_address"`
}

type Transaction struct {
	BlockHash        string   `json:"block_hash"`
	BlockNumber      string   `json:"block_number"`
	FromAccount      *Account `json:"from_account"`
	Id               int64    `json:"id"`
	Timestamp        string   `json:"timestamp"`
	ToAccount        *Account `json:"to_account"`
	TransactionHash  string   `json:"transaction_hash"`
	TransactionIndex string   `json:"transaction_index"`
}

type NftOrder struct {
	Asset                *NftAsset     `json:"asset" bson:"asset"`
	AssetBundle          *Bundle       `json:"asset_bundle" bson:"asset_bundle"`
	CreatedDate          string        `json:"created_date" bson:"created_date"`
	ClosingDate          string        `json:"closing_date" bson:"closing_date"`
	ClosingExtendable    bool          `json:"closing_extendable" bson:"closing_extendable"`
	ExpirationTime       int           `json:"expiration_time" bson:"expiration_time"`
	ListingTime          int           `json:"listing_time" bson:"listing_time"`
	OrderHash            string        `json:"order_hash" bson:"order_hash"`
	Metadata             interface{}   `json:"metadata" bson:"metadata"`
	Exchange             string        `json:"exchange" bson:"exchange"`
	Maker                *Account      `json:"maker" bson:"maker"`
	Taker                *Account      `json:"taker" bson:"taker"`
	CurrentPrice         string        `json:"current_price" bson:"current_price"`
	CurrentBounty        string        `json:"current_bounty" bson:"current_bounty"`
	BountyMultiple       string        `json:"bounty_multiple" bson:"bounty_multiple"`
	MakerRelayerFee      string        `json:"maker_relayer_fee" bson:"maker_relayer_fee"`
	TakerRelayerFee      string        `json:"taker_relayer_fee" bson:"taker_relayer_fee"`
	MakerProtocolFee     string        `json:"maker_protocol_fee" bson:"maker_protocol_fee"`
	TakerProtocolFee     string        `json:"taker_protocol_fee" bson:"taker_protocol_fee"`
	MakerReferrerFee     string        `json:"maker_referrer_fee" bson:"maker_referrer_fee"`
	FeeRecipient         *Account      `json:"fee_recipient" bson:"fee_recipient"`
	FeeMethod            int           `json:"fee_method" bson:"fee_method"`
	Side                 int           `json:"side" bson:"side"`
	SaleKind             int           `json:"sale_kind" bson:"sale_kind"`
	Target               string        `json:"target" bson:"target"`
	HowToCall            int           `json:"how_to_call" bson:"how_to_call"`
	CallData             string        `json:"calldata" bson:"call_data"`
	ReplacementPattern   string        `json:"replacement_pattern" bson:"replacement_pattern"`
	StaticTarget         string        `json:"static_target" bson:"static_target"`
	StaticExtradata      string        `json:"static_extradata" bson:"static_extradata"`
	PaymentToken         string        `json:"payment_token" bson:"payment_token"`
	PaymentTokenContract *PaymentToken `json:"payment_token_contract" bson:"payment_token_contract"`
	BasePrice            string        `json:"base_price" bson:"base_price"`
	Extra                string        `json:"extra" bson:"extra"`
	Quantity             string        `json:"quantity" bson:"quantity"`
	Salt                 string        `json:"salt" bson:"salt"`
	V                    int           `json:"v" bson:"v"`
	R                    string        `json:"r" bson:"r"`
	S                    string        `json:"s" bson:"s"`
	ApprovedOnChain      bool          `json:"approved_on_chain" bson:"approved_on_chain"`
	Cancelled            bool          `json:"cancelled" bson:"cancelled"`
	Finalized            bool          `json:"finalized" bson:"finalized"`
	MarkedInvalid        bool          `json:"marked_invalid" bson:"marked_invalid"`
	PrefixedHash         string        `json:"prefixed_hash" bson:"prefixed_hash"`
}
