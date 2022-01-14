package message

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeStatusMeta struct {
	UID             uint64             `bson:"uid,omitempty"`
	StatusID        primitive.ObjectID `bson:"status_id,omitempty"`
	StatusContent   string             `bson:"status_content,omitempty"`
	StatusImagePath string             `bson:"status_image_path,omitempty"`
}
