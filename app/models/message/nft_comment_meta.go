package message

import "go.mongodb.org/mongo-driver/bson/primitive"

type NftAssetCommentMeta struct {
	UID            uint64             `bson:"uid,omitempty"`
	GroupID        primitive.ObjectID `bson:"group_id,omitempty"`
	CommentID      primitive.ObjectID `bson:"comment_id,omitempty"`
	Content        string             `bson:"content,omitempty"`
	ParentContent  string             `bson:"parent_content,omitempty"`
	ParentUsername string             `bson:"parent_username,omitempty"`
	NftAssetName   string             `bson:"nft_asset_name,omitempty"`
	NftAssetImage  string             `bson:"nft_asset_image,omitempty"`
}
