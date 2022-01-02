package message

import "go.mongodb.org/mongo-driver/bson/primitive"

type ForwardMeta struct {
	UID      uint64             `bson:"uid,omitempty"`
	StatusID primitive.ObjectID `bson:"status_id,omitempty"`
	Content  string             `bson:"content,omitempty"`
}
