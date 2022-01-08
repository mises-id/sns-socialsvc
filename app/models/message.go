package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/message"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListMessageParams struct {
	UID        uint64
	PageParams *pagination.PageQuickParams
}

type ReadMessageParams struct {
	UID        uint64
	MessageIDs []primitive.ObjectID
	LatestID   primitive.ObjectID
}

type CreateMessageParams struct {
	UID         uint64
	MessageType enum.MessageType
	MetaData    *message.MetaData
}

type Message struct {
	message.MetaData
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UID         uint64             `bson:"uid,omitempty"`
	MessageType enum.MessageType   `bson:"message_type,omitempty"`
	ReadTime    *time.Time         `bson:"read_time,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}

func (m *Message) State() string {
	if m.ReadTime == nil {
		return "unread"
	}
	return "read"
}

func (m *Message) BeforeCreate(ctx context.Context) error {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

func CreateMessage(ctx context.Context, params *CreateMessageParams) (*Message, error) {
	message := &Message{
		UID:         params.UID,
		MessageType: params.MessageType,
		MetaData:    *params.MetaData,
	}
	var err error
	if err = message.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	if err = db.ODM(ctx).Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func ListMessage(ctx context.Context, params *ListMessageParams) ([]*Message, pagination.Pagination, error) {
	if params.PageParams == nil {
		params.PageParams = pagination.DefaultQuickParams()
	}
	messages := make([]*Message, 0)
	chain := db.ODM(ctx).Where(bson.M{"uid": params.UID})
	paginator := pagination.NewQuickPaginator(params.PageParams.Limit, params.PageParams.NextID, chain)
	page, err := paginator.Paginate(&messages)
	if err != nil {
		return nil, nil, err
	}
	return messages, page, preloadMessageData(ctx, messages...)
}

func ReadMessages(ctx context.Context, params *ReadMessageParams) error {
	query := bson.M{"uid": params.UID}
	if params.LatestID != primitive.NilObjectID {
		query["_id"] = bson.M{"$lte": params.LatestID}
	} else if params.MessageIDs != nil && len(primitive.NilObjectID) > 0 {
		query["_id"] = bson.M{"$in": params.MessageIDs}
	} else {
		return nil
	}
	_, err := db.DB().Collection("messages").UpdateMany(ctx, query, bson.D{
		{Key: "$set", Value: bson.D{{Key: "read_time", Value: time.Now()}}}})
	return err
}

func preloadMessageData(ctx context.Context, messages ...*Message) error {
	return nil
}
