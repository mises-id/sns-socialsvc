package models

import (
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	TagableID   string             `bson:"tagable_id"`
	TagableType string             `bson:"tagable_type"`
	TagType     enum.TagType       `bson:"tag_type"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}
