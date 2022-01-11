package status

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/meta"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateStatusParams struct {
	StatusType   string
	ParentID     primitive.ObjectID
	Content      string
	IsPrivate    bool
	ShowDuration int64
	Meta         meta.MetaData
	FromType     enum.FromType
}

type ListStatusParams struct {
	*pagination.PageQuickParams
	CurrentUID uint64
	UID        uint64
	ParentID   primitive.ObjectID
	FromTypes  []enum.FromType
}

func GetStatus(ctx context.Context, currentUID uint64, id primitive.ObjectID) (*models.Status, error) {
	ctxWithUID := context.WithValue(ctx, utils.CurrentUIDKey{}, currentUID)
	status, err := models.FindStatus(ctxWithUID, id)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func ListStatus(ctx context.Context, params *ListStatusParams) ([]*models.Status, pagination.Pagination, error) {
	ctxWithUID := context.WithValue(ctx, utils.CurrentUIDKey{}, params.CurrentUID)

	uids := make([]uint64, 0)
	if params.UID != 0 {
		uids = append(uids, params.UID)
		err := models.MarkFollowRead(ctx, params.CurrentUID, params.UID)
		if err != nil {
			return nil, nil, err
		}
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
	return statues, page, nil
}

func UserTimeline(ctx context.Context, uid uint64, pageParams *pagination.PageQuickParams) ([]*models.Status, pagination.Pagination, error) {
	ctxWithUID := context.WithValue(ctx, utils.CurrentUIDKey{}, uid)
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
		OnlyShow:       true,
	})
	if err != nil {
		return nil, nil, err
	}
	return statues, page, nil
}

func RecommendStatus(ctx context.Context, uid uint64, pageParams *pagination.PageQuickParams) ([]*models.Status, pagination.Pagination, error) {
	ctxWithUID := context.WithValue(ctx, utils.CurrentUIDKey{}, uid)
	statues, page, err := models.ListStatus(ctxWithUID, &models.ListStatusParams{
		UIDs:           nil,
		ParentStatusID: primitive.NilObjectID,
		FromTypes:      []enum.FromType{enum.FromPost, enum.FromForward},
		PageParams:     pageParams,
		OnlyShow:       true,
	})
	if err != nil {
		return nil, nil, err
	}
	return statues, page, nil
}

func CreateStatus(ctx context.Context, uid uint64, params *CreateStatusParams) (*models.Status, error) {
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	statusType, err := enum.StatusTypeFromString(params.StatusType)
	if err != nil {
		return nil, err
	}
	status, err := models.CreateStatus(ctx, &models.CreateStatusParams{
		UID:        uid,
		StatusType: statusType,
		Content:    params.Content,
		ParentID:   params.ParentID,
		FromType:   params.FromType,
		MetaData:   params.Meta,
	})
	if err != nil {
		return nil, err
	}
	// update user latest post time
	if err = user.UpdatePostTime(ctx, status.CreatedAt); err != nil {
		return nil, err
	}
	return status, nil
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

type ListLikeStatusParams struct {
	UID        uint64
	PageParams *pagination.PageQuickParams
}

func ListLikeStatus(ctx context.Context, params *ListLikeStatusParams) ([]*models.Like, pagination.Pagination, error) {
	return models.ListLike(ctx, params.UID, enum.LikeStatus, params.PageParams)
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
