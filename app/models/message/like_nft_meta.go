package message

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeNftAssetMeta struct {
	UID           uint64             `bson:"uid,omitempty"`
	NftAssetID    primitive.ObjectID `bson:"nft_asset_id,omitempty"`
	NftAssetName  string             `bson:"nft_asset_name,omitempty"`
	NftAssetImage string             `bson:"nft_asset_image,omitempty"`
}
