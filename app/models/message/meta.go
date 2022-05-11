package message

type MetaData struct {
	CommentMeta        *CommentMeta             `bson:"comment_meta,omitempty"`
	LikeStatusMeta     *LikeStatusMeta          `bson:"like_status_meta,omitempty"`
	LikeCommentMeta    *LikeCommentMeta         `bson:"like_comment_meta,omitempty"`
	NftCommentMeta     *NftAssetCommentMeta     `bson:"nft_comment_meta,omitempty"`
	LikeNftMeta        *LikeNftAssetMeta        `bson:"like_nft_meta,omitempty"`
	LikeNftCommentMeta *LikeNftAssetCommentMeta `bson:"like_nft_comment_meta,omitempty"`
	FansMeta           *FansMeta                `bson:"fans_meta,omitempty"`
	ForwardMeta        *ForwardMeta             `bson:"forward_meta,omitempty"`
}
