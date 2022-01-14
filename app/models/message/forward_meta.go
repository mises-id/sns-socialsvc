package message

import "go.mongodb.org/mongo-driver/bson/primitive"

type ForwardMeta struct {
	UID            uint64             `bson:"uid,omitempty"`
	StatusID       primitive.ObjectID `bson:"status_id,omitempty"`
	ForwardContent string             `bson:"forward_content,omitempty"`
	ContentSummary string             `bson:"content_summary,omitempty"`
	ImagePath      string             `bson:"image_path,omitempty"`
}
