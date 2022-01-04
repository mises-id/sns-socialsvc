package message

import (
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LikeMeta struct {
	UID        uint64              `bson:"uid,omitempty"`
	TargetID   primitive.ObjectID  `bson:"target_id,omitempty"`
	TargetType enum.LikeTargetType `bson:"target_type"`
}
