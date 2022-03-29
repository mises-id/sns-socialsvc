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

type Like struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty"`
	OwnerID    uint64              `bson:"owner_id,omitempty"`
	UID        uint64              `bson:"uid,omitempty"`
	TargetID   primitive.ObjectID  `bson:"target_id,omitempty"`
	TargetType enum.LikeTargetType `bson:"target_type"`
	DeletedAt  time.Time           `bson:"deleted_at,omitempty"`
	CreatedAt  time.Time           `bson:"created_at,omitempty"`
	UpdatedAt  time.Time           `bson:"updated_at,omitempty"`
	Status     *Status             `bson:"-"`
	Comment    *Comment            `bson:"-"`
}

func (l *Like) AfterCreate(ctx context.Context) error {
	err := l.incrUserCounter(ctx)
	if err != nil {
		return err
	}
	if err = PreloadLikeData(ctx, l); err != nil {
		return err
	}
	err = l.notifyLikeUser(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (l *Like) notifyLikeUser(ctx context.Context) error {
	err := PreloadLikeData(ctx, l)
	if err != nil {
		return err
	}
	metaData := &message.MetaData{}
	messageType := enum.NewLikeStatus
	var statusID primitive.ObjectID
	if l.TargetType == enum.LikeStatus {
		statusID = l.TargetID
		metaData.LikeStatusMeta = &message.LikeStatusMeta{
			UID:             l.UID,
			StatusID:        l.TargetID,
			StatusContent:   l.Status.ContentSummary(),
			StatusImagePath: l.Status.FirstImage(),
		}
	} else if l.TargetType == enum.LikeComment {
		statusID = l.Comment.StatusID
		metaData.LikeCommentMeta = &message.LikeCommentMeta{
			UID:             l.UID,
			CommentID:       l.TargetID,
			CommentUsername: l.Comment.User.Username,
			CommentContent:  l.Comment.Content,
		}
		messageType = enum.NewLikeComment
	}
	_, err = CreateMessage(ctx, &CreateMessageParams{
		UID:         l.OwnerID,
		StatusID:    statusID,
		FromUID:     l.UID,
		MessageType: messageType,
		MetaData:    metaData,
	})
	if err != nil {
		return err
	}
	return nil
}

func (l *Like) incrUserCounter(ctx context.Context) error {
	result := db.DB().Collection("users").FindOneAndUpdate(ctx, bson.M{"_id": l.OwnerID},
		bson.D{{
			Key: "$inc",
			Value: bson.D{{
				Key:   "liked_count",
				Value: 1,
			}}},
		})
	return result.Err()
}

func CreateLike(ctx context.Context, uid, ownerID uint64, targetID primitive.ObjectID, targetType enum.LikeTargetType) (*Like, error) {
	like := &Like{
		OwnerID:    ownerID,
		UID:        uid,
		TargetID:   targetID,
		TargetType: targetType,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := db.ODM(ctx).Create(like).Error
	if err != nil {
		return nil, err
	}
	return like, like.AfterCreate(ctx)
}

func DeleteLike(ctx context.Context, id primitive.ObjectID) error {
	return db.DB().Collection("likes").FindOneAndUpdate(ctx, bson.M{
		"_id":        id,
		"deleted_at": bson.M{"$exists": false},
	}, bson.M{"$set": bson.M{"deleted_at": time.Now()}}).Err()
}

func FindLike(ctx context.Context, uid uint64, targetID primitive.ObjectID, targetType enum.LikeTargetType) (*Like, error) {
	like := &Like{}
	err := db.ODM(ctx).Where(bson.M{
		"uid":         uid,
		"target_id":   targetID,
		"target_type": targetType,
		"deleted_at":  bson.M{"$exists": false},
	}).First(like).Error
	if err != nil {
		return nil, err
	}
	if err = PreloadLikeData(ctx, like); err != nil {
		return nil, err
	}
	return like, err
}

func GetLikeMap(ctx context.Context, uid uint64, targetIDs []primitive.ObjectID, targetType enum.LikeTargetType, preloadData bool) (map[primitive.ObjectID]*Like, error) {
	likes := make([]*Like, 0)
	err := db.ODM(ctx).Where(bson.M{
		"uid":         uid,
		"target_id":   bson.M{"$in": targetIDs},
		"target_type": targetType,
		"deleted_at":  nil,
	}).Find(&likes).Error
	if err != nil {
		return nil, err
	}
	if preloadData {
		if err = PreloadLikeData(ctx, likes...); err != nil {
			return nil, err
		}
	}
	likeMap := make(map[primitive.ObjectID]*Like)
	for _, like := range likes {
		likeMap[like.TargetID] = like
	}
	return likeMap, nil
}

func ListLike(ctx context.Context, uid uint64, tp enum.LikeTargetType, pageParams *pagination.PageQuickParams) ([]*Like, pagination.Pagination, error) {
	if pageParams == nil {
		pageParams = pagination.DefaultQuickParams()
	}
	likes := make([]*Like, 0)
	chain := db.ODM(ctx).Where(bson.M{"uid": uid, "target_type": tp, "deleted_at": nil})
	blockedUIDs, err := ListBlockedUIDs(ctx)
	if err != nil {
		return nil, nil, err
	}
	if len(blockedUIDs) > 0 {
		chain = chain.Where(bson.M{"target_id": bson.M{"$nin": blockedUIDs}})
	}
	paginator := pagination.NewQuickPaginator(pageParams.Limit, pageParams.NextID, chain)
	page, err := paginator.Paginate(&likes)
	if err != nil {
		return nil, nil, err
	}
	return likes, page, PreloadLikeData(ctx, likes...)
}

func PreloadLikeData(ctx context.Context, likes ...*Like) error {
	err := preloadLikeComment(ctx, likes...)
	if err != nil {
		return err
	}
	err = preloadLikeStatus(ctx, likes...)
	if err != nil {
		return err
	}
	return nil
}

func preloadLikeStatus(ctx context.Context, likes ...*Like) error {

	var uid uint64
	currentUID, ok := ctx.Value(utils.CurrentUIDKey{}).(uint64)
	if ok {
		uid = currentUID
	}
	statusIDs := make([]primitive.ObjectID, 0)
	for _, like := range likes {
		if like.TargetType != enum.LikeStatus {
			continue
		}
		statusIDs = append(statusIDs, like.TargetID)
	}
	statuses, err := FindStatusByIDs(ctx, statusIDs...)
	if err != nil {
		return err
	}
	statusMap := make(map[primitive.ObjectID]*Status)
	for _, status := range statuses {
		statusMap[status.ID] = status
	}
	for _, like := range likes {
		if like.TargetType != enum.LikeStatus {
			continue
		}
		if statusMap[like.TargetID] != nil && (statusMap[like.TargetID].IsPublic || uid == statusMap[like.TargetID].UID) {
			like.Status = statusMap[like.TargetID]
		}

	}
	return nil
}

func preloadLikeComment(ctx context.Context, likes ...*Like) error {
	commentIDs := make([]primitive.ObjectID, 0)
	for _, like := range likes {
		if like.TargetType != enum.LikeComment {
			continue
		}
		commentIDs = append(commentIDs, like.TargetID)
	}
	comments, err := FindCommentByIDs(ctx, commentIDs...)
	if err != nil {
		return err
	}
	if err = PreloadCommentData(ctx, comments...); err != nil {
		return err
	}
	commentMap := make(map[primitive.ObjectID]*Comment)
	for _, comment := range comments {
		commentMap[comment.ID] = comment
	}
	for _, like := range likes {
		if like.TargetType != enum.LikeComment {
			continue
		}
		like.Comment = commentMap[like.TargetID]
	}
	return nil
}
