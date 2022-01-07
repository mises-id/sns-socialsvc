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

type ListCommentParams struct {
	StatusID   primitive.ObjectID
	GroupID    primitive.ObjectID
	OpponentID uint64
	TargetUID  uint64
	UID        uint64
	PageParams *pagination.PageQuickParams
}

type CreateCommentParams struct {
	StatusID   primitive.ObjectID
	ParentID   primitive.ObjectID
	GroupID    primitive.ObjectID
	OpponentID uint64
	UID        uint64
	Content    string
}

type Comment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	StatusID   primitive.ObjectID `bson:"status_id,omitempty"`
	ParentID   primitive.ObjectID `bson:"parent_id,omitempty"`
	GroupID    primitive.ObjectID `bson:"group_id,omitempty"`
	OpponentID uint64             `bson:"opponent_uid,omitempty"`
	UID        uint64             `bson:"uid,omitempty"`
	Content    string             `bson:"content,omitempty"`
	CreatedAt  time.Time          `bson:"created_at,omitempty"`
	UpdatedAt  time.Time          `bson:"updated_at,omitempty"`
	User       *User              `bson:"-"`
}

func (c *Comment) BeforeCreate(ctx context.Context) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Comment) AfterCreate(ctx context.Context) error {
	_, err := CreateMessage(ctx, &CreateMessageParams{
		UID:         c.OpponentID,
		MessageType: enum.NewComment,
		MetaData: &message.MetaData{
			CommentMeta: &message.CommentMeta{
				UID:       c.UID,
				GroupID:   c.GroupID,
				CommentID: c.ID,
				Content:   c.Content,
			},
		},
	})
	return err
}

func FindComment(ctx context.Context, id primitive.ObjectID) (*Comment, error) {
	comment := &Comment{}
	err := db.ODM(ctx).First(comment, bson.M{"_id": id}).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func ListComment(ctx context.Context, params *ListCommentParams) ([]*Comment, pagination.Pagination, error) {
	if params.PageParams == nil {
		params.PageParams = pagination.DefaultQuickParams()
	}
	comments := make([]*Comment, 0)
	chain := db.ODM(ctx)
	if params.StatusID != primitive.NilObjectID {
		chain = chain.Where(bson.M{"status_id": params.StatusID})
	}
	if params.UID != 0 {
		chain = chain.Where(bson.M{"uid": params.UID})
	}
	if params.OpponentID != 0 {
		chain = chain.Where(bson.M{"opponent_uid": params.OpponentID})
	}
	if params.GroupID != primitive.NilObjectID {
		chain = chain.Where(bson.M{"group_id": params.GroupID})
	}
	paginator := pagination.NewQuickPaginator(params.PageParams.Limit, params.PageParams.NextID, chain)
	page, err := paginator.Paginate(&comments)
	if err != nil {
		return nil, nil, err
	}
	return comments, page, preloadCommentData(ctx, comments...)
}

func CreateComment(ctx context.Context, params *CreateCommentParams) (*Comment, error) {
	comment := &Comment{
		UID:        params.UID,
		StatusID:   params.StatusID,
		ParentID:   params.ParentID,
		GroupID:    params.GroupID,
		OpponentID: params.OpponentID,
		Content:    params.Content,
	}
	var err error
	if err = comment.BeforeCreate(ctx); err != nil {
		return nil, err
	}
	if err = db.ODM(ctx).Create(comment).Error; err != nil {
		return nil, err
	}
	if err = comment.AfterCreate(ctx); err != nil {
		return nil, err
	}
	return comment, preloadCommentData(ctx, comment)

}

func preloadCommentData(ctx context.Context, comments ...*Comment) error {
	userIDs := make([]uint64, len(comments))
	for i, comment := range comments {
		userIDs[i] = comment.UID
	}
	users, err := GetUserMap(ctx, userIDs...)
	if err != nil {
		return err
	}
	for _, comment := range comments {
		comment.User = users[comment.UID]
	}
	return nil
}
