package message

import "go.mongodb.org/mongo-driver/bson/primitive"

type CommentMeta struct {
	UID       uint64             `bson:"uid,omitempty"`
	GroupID   primitive.ObjectID `bson:"group_id,omitempty"`
	CommentID primitive.ObjectID `bson:"comment_id,omitempty"`
	Content   string             `bson:"content,omitempty"`
}
