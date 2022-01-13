package models

import (
	"context"
	"time"

	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/message"
	"github.com/mises-id/sns-socialsvc/lib/db"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
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
	Status     *Status
}

type Comment struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	StatusID      primitive.ObjectID   `bson:"status_id,omitempty"`
	ParentID      primitive.ObjectID   `bson:"parent_id,omitempty"`
	GroupID       primitive.ObjectID   `bson:"group_id,omitempty"`
	OpponentID    uint64               `bson:"opponent_uid,omitempty"`
	CommentIDs    []primitive.ObjectID `bson:"comment_ids,omitempty"`
	UID           uint64               `bson:"uid,omitempty"`
	LikesCount    uint64               `bson:"likes_count,omitempty"`
	CommentsCount uint64               `bson:"comments_count,omitempty"`
	Content       string               `bson:"content,omitempty"`
	CreatedAt     time.Time            `bson:"created_at,omitempty"`
	UpdatedAt     time.Time            `bson:"updated_at,omitempty"`
	User          *User                `bson:"-"`
	Opponent      *User                `bson:"-"`
	Comments      []*Comment           `bson:"-"`
	Status        *Status              `bson:"-"`
	IsLiked       bool                 `bson:"-"`
	Parent        *Comment             `bson:"-"`
}

func (c *Comment) BeforeCreate(ctx context.Context) error {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Comment) AfterCreate(ctx context.Context) error {
	var err error
	if c.ParentID.IsZero() {
		err = c.notifyStatusUser(ctx)
	} else {
		err = c.notifyCommentUser(ctx)
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) ParentContent() string {
	if c.ParentID.IsZero() {
		return ""
	}
	return c.Parent.Content
}

func (c *Comment) ParentUserName() string {
	if c.ParentID.IsZero() {
		return ""
	}
	return c.Parent.User.Username
}

func (c *Comment) notifyStatusUser(ctx context.Context) error {
	_, err := CreateMessage(ctx, &CreateMessageParams{
		UID:         c.Status.UID,
		FromUID:     c.UID,
		MessageType: enum.NewComment,
		MetaData: &message.MetaData{
			CommentMeta: &message.CommentMeta{
				UID:                  c.UID,
				GroupID:              c.GroupID,
				CommentID:            c.ID,
				Content:              c.Content,
				StatusContentSummary: c.Status.ContentSummary(),
				StatusImagePath:      c.Status.FirstImage(),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) notifyCommentUser(ctx context.Context) error {
	if c.UID == c.OpponentID {
		return nil
	}
	var err error
	c.Parent, err = FindComment(ctx, c.ParentID)
	if err != nil {
		return err
	}
	_, err = CreateMessage(ctx, &CreateMessageParams{ // notify parent comment user
		UID:         c.OpponentID,
		FromUID:     c.UID,
		MessageType: enum.NewComment,
		MetaData: &message.MetaData{
			CommentMeta: &message.CommentMeta{
				UID:                  c.UID,
				GroupID:              c.GroupID,
				CommentID:            c.ID,
				Content:              c.Content,
				ParentContent:        c.ParentContent(),
				ParentUserName:       c.ParentUserName(),
				StatusContentSummary: c.Status.ContentSummary(),
				StatusImagePath:      c.Status.FirstImage(),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) IncCommentCounter(ctx context.Context, counterKey string, values ...int) error {
	if counterKey == "" {
		return nil
	}
	value := 1
	if len(values) > 0 {
		value = values[0]
	}
	result := db.DB().Collection("comments").FindOneAndUpdate(ctx, bson.M{"_id": c.ID},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   counterKey,
				Value: value,
			}}},
		})
	if err := result.Err(); err != nil {
		return err
	}
	return result.Decode(c)
}

func (c *Comment) AddChildComment(ctx context.Context, commentID primitive.ObjectID) error {
	result := db.DB().Collection("comments").FindOneAndUpdate(ctx, bson.M{"_id": c.ID},
		bson.D{{
			Key: "$push",
			Value: bson.D{{
				Key:   "comment_ids",
				Value: commentID,
			}}},
		})
	if err := result.Err(); err != nil {
		return err
	}
	return result.Decode(c)
}

func FindComment(ctx context.Context, id primitive.ObjectID) (*Comment, error) {
	comment := &Comment{}
	err := db.ODM(ctx).First(comment, bson.M{"_id": id}).Error
	if err != nil {
		return nil, err
	}
	return comment, PreloadCommentData(ctx, comment)
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
	} else {
		chain = chain.Where(bson.M{"group_id": bson.M{"$exists": false}})
	}
	chain = chain.Sort(bson.M{"_id": 1})
	paginator := pagination.NewQuickPaginator(params.PageParams.Limit, params.PageParams.NextID, chain)
	page, err := paginator.Paginate(&comments)
	if err != nil {
		return nil, nil, err
	}
	return comments, page, PreloadCommentData(ctx, comments...)
}

func CreateComment(ctx context.Context, params *CreateCommentParams) (*Comment, error) {
	comment := &Comment{
		UID:        params.UID,
		StatusID:   params.StatusID,
		ParentID:   params.ParentID,
		GroupID:    params.GroupID,
		OpponentID: params.OpponentID,
		Content:    params.Content,
		Status:     params.Status,
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
	return comment, PreloadCommentData(ctx, comment)

}

func PreloadCommentData(ctx context.Context, comments ...*Comment) error {
	if err := preloadCommentUser(ctx, comments...); err != nil {
		return err
	}
	if err := preloadCommentLikeState(ctx, comments...); err != nil {
		return err
	}
	return preloadCommentChildren(ctx, comments...)
}

func preloadCommentLikeState(ctx context.Context, comments ...*Comment) error {
	currentUID, ok := ctx.Value(utils.CurrentUIDKey{}).(uint64)
	if !ok || currentUID == 0 {
		return nil
	}
	ids := make([]primitive.ObjectID, len(comments))
	for i, comment := range comments {
		ids[i] = comment.ID
	}
	likeMap, err := GetLikeMap(ctx, currentUID, ids, enum.LikeComment, true)
	if err != nil {
		return err
	}
	for _, comment := range comments {
		comment.IsLiked = likeMap[comment.ID] != nil
	}
	return nil
}

func preloadCommentChildren(ctx context.Context, comments ...*Comment) error {
	ids := make([]primitive.ObjectID, 0)
	for _, comment := range comments {
		if comment.CommentIDs == nil {
			continue
		}
		ids = append(ids, comment.CommentIDs...)
	}
	children, err := FindCommentByIDs(ctx, ids...)
	if err != nil {
		return err
	}
	if err = preloadCommentUser(ctx, children...); err != nil {
		return err
	}
	childrenMap := make(map[primitive.ObjectID]*Comment)
	for _, child := range children {
		childrenMap[child.ID] = child
	}
	for _, comment := range comments {
		if comment.CommentIDs == nil {
			continue
		}
		comment.Comments = make([]*Comment, 0)
		for _, id := range comment.CommentIDs {
			if childrenMap[id] != nil {
				comment.Comments = append(comment.Comments, childrenMap[id])
			}
		}
	}
	return nil
}

func preloadCommentUser(ctx context.Context, comments ...*Comment) error {
	userIDs := make([]uint64, 0)
	for _, comment := range comments {
		userIDs = append(userIDs, comment.UID)
		if comment.OpponentID != 0 {
			userIDs = append(userIDs, comment.OpponentID)
		}
	}
	users, err := GetUserMap(ctx, userIDs...)
	if err != nil {
		return err
	}
	for _, comment := range comments {
		comment.User = users[comment.UID]
		if !comment.GroupID.IsZero() && comment.ParentID != comment.GroupID {
			comment.Opponent = users[comment.OpponentID]
		}
	}
	return nil
}

func FindCommentByIDs(ctx context.Context, ids ...primitive.ObjectID) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	return comments, db.ODM(ctx).Where(bson.M{"_id": bson.M{"$in": ids}}).Find(&comments).Error
}
