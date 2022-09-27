package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/lib/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CreateComplaintParams struct {
		UID        uint64 `json:"uid"`
		TargetType string `json:"target_type"`
		TargetID   string `json:"target_id"`
		Reason     string `json:"reason"`
	}

	Complaint struct {
		ID         primitive.ObjectID `bson:"_id,omitempty"`
		UID        uint64             `bson:"uid"`
		TargetType string             `bson:"target_type"`
		TargetID   string             `bson:"target_id"`
		Reason     string             `bson:"reason"`
		CreatedAt  time.Time          `bson:"created_at,omitempty"`
	}
)

func CreateComplaint(ctx context.Context, data *CreateComplaintParams) (*Complaint, error) {
	insert := &Complaint{
		UID:        data.UID,
		TargetType: data.TargetType,
		TargetID:   data.TargetID,
		Reason:     data.Reason,
	}
	err := db.ODM(ctx).Create(insert).Error
	if err != nil {
		return nil, err
	}
	return insert, err
}
