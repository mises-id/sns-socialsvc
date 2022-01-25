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
	"go.mongodb.org/mongo-driver/mongo"
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
	FromUID     uint64
	StatusID    primitive.ObjectID
	MessageType enum.MessageType
	MetaData    *message.MetaData
}

type Message struct {
	message.MetaData
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	UID              uint64             `bson:"uid,omitempty"`
	FromUID          uint64             `bson:"from_uid,omitempty"`
	StatusID         primitive.ObjectID `bson:"status_id,omitempty"`
	MessageType      enum.MessageType   `bson:"message_type,omitempty"`
	ReadTime         *time.Time         `bson:"read_time"`
	CreatedAt        time.Time          `bson:"created_at,omitempty"`
	UpdatedAt        time.Time          `bson:"updated_at,omitempty"`
	FromUser         *User              `bson:"-"`
	Status           *Status            `bson:"-"`
	Comment          *Comment           `bson:"-"`
	StatusIsDeleted  bool               `bson:"-"`
	CommentIsDeleted bool               `bson:"-"`
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
	if params.UID == params.FromUID {
		return nil, nil
	}
	message := &Message{
		UID:         params.UID,
		StatusID:    params.StatusID,
		FromUID:     params.FromUID,
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
	return messages, page, PreloadMessageData(ctx, messages...)
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

func UnreadMessagesCount(ctx context.Context, uid uint64) (uint32, error) {
	var c int64
	return uint32(c), db.ODM(ctx).Model(&Message{}).Where(bson.M{"uid": uid, "read_time": nil}).Count(&c).Error
}

func LatestUnreadMessage(ctx context.Context, uid uint64) (*Message, error) {
	message := &Message{}
	err := db.ODM(ctx).Model(&Message{}).
		Where(bson.M{"uid": uid, "read_time": nil}).Sort(bson.M{"_id": -1}).First(message).Error
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return message, PreloadMessageData(ctx, message)
}

func PreloadMessageData(ctx context.Context, messages ...*Message) error {
	if err := preloadMessageUser(ctx, messages...); err != nil {
		return err
	}
	if err := preloadMessageStatus(ctx, messages...); err != nil {
		return err
	}
	if err := preloadMessageComment(ctx, messages...); err != nil {
		return err
	}
	return nil
}

func preloadMessageUser(ctx context.Context, messages ...*Message) error {
	userIDs := make([]uint64, len(messages))
	for i, message := range messages {
		userIDs[i] = message.FromUID
	}
	users, err := GetUserMap(ctx, userIDs...)
	if err != nil {
		return err
	}
	for _, message := range messages {
		message.FromUser = users[message.FromUID]
	}
	return nil
}

func preloadMessageStatus(ctx context.Context, messages ...*Message) error {
	statusIDs := make([]primitive.ObjectID, 0)
	for _, message := range messages {
		if !message.StatusID.IsZero() {
			statusIDs = append(statusIDs, message.StatusID)
		}
	}
	statuses, err := FindStatusByIDs(ctx, statusIDs...)
	if err != nil {
		return err
	}
	statusMap := map[primitive.ObjectID]*Status{}
	for _, status := range statuses {
		statusMap[status.ID] = status
	}
	for _, message := range messages {
		if !message.StatusID.IsZero() {
			message.Status = statusMap[message.StatusID]
			if message.Status == nil {
				message.StatusIsDeleted = true
			}
		}
	}
	return nil
}
func preloadMessageComment(ctx context.Context, messages ...*Message) error {
	commentIDs := make([]primitive.ObjectID, 0)
	for _, message := range messages {
		if message.CommentMeta != nil && !message.CommentMeta.CommentID.IsZero() {
			commentIDs = append(commentIDs, message.CommentMeta.CommentID)
		}
	}
	comments, err := FindCommentByIDs(ctx, commentIDs...)
	if err != nil {
		return err
	}
	commentMap := map[primitive.ObjectID]*Comment{}
	for _, comment := range comments {
		commentMap[comment.ID] = comment
	}
	for _, message := range messages {
		if message.CommentMeta != nil && !message.CommentMeta.CommentID.IsZero() {
			//message.Comment = commentMap[message.CommentMeta.CommentID]
			if commentMap[message.CommentMeta.CommentID] == nil {
				message.CommentIsDeleted = true
			}
		}
	}
	return nil
}
