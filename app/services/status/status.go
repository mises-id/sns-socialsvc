package status

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/meta"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateStatusParams struct {
	StatusType string
	ParentID   primitive.ObjectID
	Content    string
	Meta       *meta.MetaData
	FromType   enum.FromType
}

type ListStatusParams struct {
	*pagination.PageQuickParams
	CurrentUID uint64
	UID        uint64
	ParentID   primitive.ObjectID
	FromTypes  []enum.FromType
}

func GetStatus(ctx context.Context, currentUID uint64, id primitive.ObjectID) (*models.Status, error) {
	ctxWithUID := context.WithValue(ctx, "CurrentUID", currentUID)
	status, err := models.FindStatus(ctxWithUID, id)
	if err != nil {
		return nil, err
	}
	return status, batchSetIsLiked(ctx, currentUID, status)
}

func ListStatus(ctx context.Context, params *ListStatusParams) ([]*models.Status, pagination.Pagination, error) {
	ctxWithUID := context.WithValue(ctx, "CurrentUID", params.CurrentUID)

	uids := make([]uint64, 0)
	if params.UID != 0 {
		uids = append(uids, params.UID)
	}
	listParams := &models.ListStatusParams{
		UIDs:           uids,
		ParentStatusID: params.ParentID,
		PageParams:     params.PageQuickParams,
		FromTypes:      params.FromTypes,
	}
	statues, page, err := models.ListStatus(ctxWithUID, listParams)
	if err != nil {
		return nil, nil, err
	}
	return statues, page, batchSetIsLiked(ctx, params.CurrentUID, statues...)
}

func UserTimeline(ctx context.Context, uid uint64, pageParams *pagination.PageQuickParams) ([]*models.Status, pagination.Pagination, error) {
	ctxWithUID := context.WithValue(ctx, "CurrentUID", uid)
	friendIDs, err := models.ListFollowingUserIDs(ctx, uid)
	if err != nil {
		return nil, nil, err
	}
	if len(friendIDs) == 0 {
		return []*models.Status{}, &pagination.QuickPagination{
			Limit: pageParams.Limit,
		}, nil
	}

	statues, page, err := models.ListStatus(ctxWithUID, &models.ListStatusParams{
		UIDs:           friendIDs,
		ParentStatusID: primitive.NilObjectID,
		PageParams:     pageParams,
	})
	if err != nil {
		return nil, nil, err
	}
	return statues, page, batchSetIsLiked(ctx, uid, statues...)
}

func RecommendStatus(ctx context.Context, uid uint64, pageParams *pagination.PageQuickParams) ([]*models.Status, pagination.Pagination, error) {
	ctxWithUID := context.WithValue(ctx, "CurrentUID", uid)
	statues, page, err := models.ListStatus(ctxWithUID, &models.ListStatusParams{
		UIDs:           nil,
		ParentStatusID: primitive.NilObjectID,
		FromTypes:      []enum.FromType{enum.FromPost, enum.FromForward},
		PageParams:     pageParams,
	})
	if err != nil {
		return nil, nil, err
	}
	return statues, page, batchSetIsLiked(ctx, uid, statues...)
}

func CreateStatus(ctx context.Context, uid uint64, params *CreateStatusParams) (*models.Status, error) {
	statusType, err := enum.StatusTypeFromString(params.StatusType)
	if err != nil {
		return nil, err
	}
	return models.CreateStatus(ctx, &models.CreateStatusParams{
		UID:        uid,
		StatusType: statusType,
		Content:    params.Content,
		ParentID:   params.ParentID,
		FromType:   params.FromType,
		MetaData:   params.Meta,
	})
}

func LikeStatus(ctx context.Context, uid uint64, statusID primitive.ObjectID) (*models.Like, error) {
	status, err := models.FindStatus(ctx, statusID)
	if err != nil {
		return nil, err
	}
	like, err := models.FindLike(ctx, uid, statusID, enum.LikeStatus)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == nil {
		return like, nil
	}
	like, err = models.CreateLike(ctx, uid, status.UID, statusID, enum.LikeStatus)
	if err != nil {
		return nil, err
	}
	return like, status.IncStatusCounter(ctx, "likes_count")
}

func UnlikeStatus(ctx context.Context, uid uint64, statusID primitive.ObjectID) error {
	like, err := models.FindLike(ctx, uid, statusID, enum.LikeStatus)
	if err != nil {
		return err
	}
	status, err := models.FindStatus(ctx, statusID)
	if err != nil {
		return err
	}
	if err = models.DeleteLike(ctx, like.ID); err != nil {
		return err
	}
	return status.IncStatusCounter(ctx, "likes_count", -1)
}

func DeleteStatus(ctx context.Context, uid uint64, id primitive.ObjectID) error {
	status, err := models.FindStatus(ctx, id)
	if err != nil {
		return err
	}
	if status.UID != uid {
		return codes.ErrForbidden
	}
	return models.DeleteStatus(ctx, id)
}

func batchSetIsLiked(ctx context.Context, uid uint64, statuses ...*models.Status) error {
	if uid == 0 {
		return nil
	}
	statusIDs := make([]primitive.ObjectID, len(statuses))
	for i, status := range statuses {
		statusIDs[i] = status.ID
	}
	likeMap, err := models.GetStatusLikeMap(ctx, uid, statusIDs)
	if err != nil {
		return err
	}
	for _, status := range statuses {
		status.IsLiked = likeMap[status.ID] != nil
	}
	return nil
}
