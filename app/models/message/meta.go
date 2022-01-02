package message

type MetaData struct {
	CommentMeta *CommentMeta `bson:"comment_meta"`
	LikeMeta    *LikeMeta    `bson:"like_meta"`
	FansMeta    *FansMeta    `bson:"fans_meta"`
	ForwardMeta *ForwardMeta `bson:"forward_meta"`
}
