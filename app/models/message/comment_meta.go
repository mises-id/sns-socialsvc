package message

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommentMeta struct {
	UID                  uint64             `bson:"uid,omitempty"`
	GroupID              primitive.ObjectID `bson:"group_id,omitempty"`
	CommentID            primitive.ObjectID `bson:"comment_id,omitempty"`
	Content              string             `bson:"content,omitempty"`
	ParentContent        string             `bson:"parent_content,omitempty"`
	ParentUsername       string             `bson:"parent_username,omitempty"`
	StatusContentSummary string             `bson:"status_content_summary,omitempty"`
	StatusImagePath      string             `bson:"status_image_path,omitempty"`
	StatusImageURL       string             `bson:"-"`
}
