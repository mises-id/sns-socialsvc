package message

type MetaData struct {
	CommentMeta *CommentMeta `bson:"comment_meta,omitempty"`
	LikeMeta    *LikeMeta    `bson:"like_meta,omitempty"`
	FansMeta    *FansMeta    `bson:"fans_meta,omitempty"`
	ForwardMeta *ForwardMeta `bson:"forward_meta,omitempty"`
}
