package follow

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"go.mongodb.org/mongo-driver/mongo"
)

func ListFriendship(ctx context.Context, uid uint64, relationType enum.RelationType, pageParams *pagination.QuickPagination) ([]*models.Follow, pagination.Pagination, error) {
	// check user exsit
	_, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, nil, err
	}
	return models.ListFollow(ctx, uid, relationType, pageParams)
}

func Follow(ctx context.Context, fromUID, toUID uint64) (*models.Follow, error) {
	if fromUID == toUID {
		return nil, codes.ErrInvalidArgument
	}
	fromUser, err := models.FindUser(ctx, fromUID)
	if err != nil {
		return nil, err
	}
	toUser, err := models.FindUser(ctx, toUID)
	if err != nil {
		return nil, err
	}
	isFriend := false
	follow, err := models.GetFollow(ctx, fromUID, toUID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	fansFollow, err := models.GetFollow(ctx, toUID, fromUID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if err == nil {
		isFriend = true
		if !fansFollow.IsFriend {
			if err = fansFollow.SetFriend(ctx, true); err != nil {
				return nil, err
			}
		}
	}
	if follow != nil {
		return follow, follow.SetFriend(ctx, isFriend)
	}
	if err = fromUser.IncFollowingCount(ctx); err != nil {
		return nil, err
	}
	if err = toUser.IncFansCount(ctx); err != nil {
		return nil, err
	}

	return models.CreateFollow(ctx, fromUID, toUID, isFriend)
}

func Unfollow(ctx context.Context, fromUID, toUID uint64) error {
	_, err := models.GetFollow(ctx, fromUID, toUID)
	if err != nil {
		return nil
	}
	fansFollow, err := models.GetFollow(ctx, toUID, fromUID)
	if err == nil && fansFollow.IsFriend {
		if err = fansFollow.SetFriend(ctx, false); err != nil {
			return err
		}
	} else if err != mongo.ErrNoDocuments {
		return err
	}
	return models.DeleteFollow(ctx, fromUID, toUID)
}
