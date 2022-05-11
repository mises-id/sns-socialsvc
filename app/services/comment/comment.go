package comment

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ListCommentParams struct {
	models.ListCommentParams
}

type CreateCommentParams struct {
	*models.CreateCommentParams
}

func GetComment(ctx context.Context, currentUID uint64, id primitive.ObjectID) (*models.Comment, error) {

	comment, err := models.FindComment(ctx, id)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func ListComment(ctx context.Context, params *ListCommentParams) ([]*models.Comment, pagination.Pagination, error) {
	return models.ListComment(ctx, &params.ListCommentParams)
}

func DeleteComment(ctx context.Context, currentUID uint64, id primitive.ObjectID) error {
	comment, err := models.FindComment(ctx, id)
	if err != nil {
		return err
	}
	if comment.UID != currentUID {
		return codes.ErrNotFound
	}
	if err = comment.Delete(ctx); err != nil {
		return err
	}
	//delete comment
	deleteCommentNum := 1
	//1.comment is comment1, delete this comment and this group comment
	if comment.ParentID.IsZero() && comment.GroupID.IsZero() {
		//delete group comment
		err := models.DeleteManyByGroupId(ctx, comment.ID)
		if err != nil {
			return err
		}
		deleteCommentNum += int(comment.CommentsCount)

	}
	if !comment.StatusID.IsZero() {
		status, err := models.FindStatus(ctx, comment.StatusID)
		if err != nil {
			return err
		}
		if err = status.IncStatusCounter(ctx, "comments_count", -deleteCommentNum); err != nil {
			return err
		}
	}
	if !comment.NftAssetID.IsZero() {
		nft_asset, err := models.FindNftAssetByID(ctx, comment.NftAssetID)
		if err != nil {
			return err
		}
		if err = nft_asset.IncNftAssetCounter(ctx, "comments_count", -deleteCommentNum); err != nil {
			return err
		}
	}
	//2.comment is comment2, delete this comment and handler group comment count
	if !comment.GroupID.IsZero() && comment.GroupID != comment.ID {
		groupComment, err := models.FindComment(ctx, comment.GroupID)
		if err != nil {
			return err
		}
		if err = groupComment.IncCommentCounter(ctx, "comments_count", -1); err != nil {
			return err
		}
		//delete comment id in group comment comment_ids and add other child comment
		var isUpdateCommentIds bool
		for _, cid := range groupComment.CommentIDs {
			if cid == id {
				isUpdateCommentIds = true
				break
			}
		}
		if isUpdateCommentIds {
			err := groupComment.RemoveChildComment(ctx, id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CreateComment(ctx context.Context, params *CreateCommentParams) (*models.Comment, error) {
	var err error
	var status *models.Status
	var nft_asset *models.NftAsset
	commentParams := params.CreateCommentParams

	if !params.NftAssetID.IsZero() {
		nft_asset, err = models.FindNftAssetByID(ctx, params.NftAssetID)
		if err != nil {
			return nil, err
		}
		commentParams.NftAsset = nft_asset
	}
	if !params.StatusID.IsZero() {
		// check status exsist
		status, err = models.FindStatus(ctx, params.StatusID)
		if err != nil {
			return nil, err
		}
		//access rights
		if status != nil && !status.IsPublic && status.UID != params.UID {
			return nil, codes.ErrForbidden
		}
		statusBlocked, err := models.IsBlocked(ctx, status.UID, params.UID)
		if err != nil {
			return nil, err
		}
		if statusBlocked {
			return nil, codes.ErrUserInBlacklist
		}
		commentParams.Status = status
	}

	var groupComment *models.Comment
	if params.ParentID != primitive.NilObjectID {
		parent, err := models.FindComment(ctx, params.ParentID)
		if err != nil {
			return nil, err
		}
		if parent.GroupID == primitive.NilObjectID {
			commentParams.GroupID = parent.ID
			groupComment = parent
		} else {
			commentParams.GroupID = parent.GroupID
			groupComment, err = models.FindComment(ctx, commentParams.GroupID)
			if err != nil {
				return nil, err
			}
		}
		commentParams.OpponentID = parent.UID
		blocked, err := models.IsBlocked(ctx, parent.UID, params.UID)
		if err != nil {
			return nil, err
		}
		if blocked {
			return nil, codes.ErrUserInBlacklist
		}
	}
	comment, err := models.CreateComment(ctx, commentParams)
	if err != nil {
		return nil, err
	}
	if err = addChildrenComment(ctx, groupComment, comment); err != nil {
		return nil, err
	}
	if err = incrCommentCounter(ctx, status, nft_asset, groupComment); err != nil {
		return nil, err
	}
	if err = models.PreloadCommentData(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func LikeComment(ctx context.Context, uid uint64, commentID primitive.ObjectID) (*models.Like, error) {
	comment, err := models.FindComment(ctx, commentID)
	if err != nil {
		return nil, err
	}
	var likeTargetType enum.LikeTargetType
	likeTargetType = enum.LikeComment
	if comment.StatusID.IsZero() {
		likeTargetType = enum.LikeNftComment
	}
	like, err := models.FindLike(ctx, uid, commentID, likeTargetType)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == nil {
		return like, nil
	}
	like, err = models.CreateLike(ctx, uid, comment.UID, commentID, likeTargetType)
	if err != nil {
		return nil, err
	}
	return like, comment.IncCommentCounter(ctx, "likes_count")
}

func UnlikeComment(ctx context.Context, uid uint64, commentID primitive.ObjectID) error {
	comment, err := models.FindComment(ctx, commentID)
	if err != nil {
		return err
	}
	var likeTargetType enum.LikeTargetType
	likeTargetType = enum.LikeComment
	if comment.StatusID.IsZero() {
		likeTargetType = enum.LikeNftComment
	}
	like, err := models.FindLike(ctx, uid, commentID, likeTargetType)
	if err != nil {
		return err
	}
	if err = models.DeleteLike(ctx, like.ID); err != nil {
		return err
	}
	return comment.IncCommentCounter(ctx, "likes_count", -1)
}

func incrCommentCounter(ctx context.Context, status *models.Status, nft_asset *models.NftAsset, comment *models.Comment) error {
	if status != nil {
		err := status.IncStatusCounter(ctx, "comments_count")
		if err != nil {
			return err
		}
	}
	if nft_asset != nil {
		err := nft_asset.IncNftAssetCounter(ctx, "comments_count")
		if err != nil {
			return err
		}
	}
	if comment != nil {
		err := comment.IncCommentCounter(ctx, "comments_count")
		if err != nil {
			return err
		}
	}
	return nil
}

func addChildrenComment(ctx context.Context, groupComment, comment *models.Comment) error {

	if groupComment == nil {
		return nil
	}
	/* if groupComment.CommentIDs != nil && len(groupComment.CommentIDs) >= 3 {
		return nil
	} */
	return groupComment.AddChildComment(ctx, comment.ID)
}
