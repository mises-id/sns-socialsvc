package user

import (
	"context"

	"github.com/mises-id/socialsvc/app/models"
	"github.com/mises-id/socialsvc/app/models/enum"
	"github.com/mises-id/socialsvc/lib/codes"
)

type UserProfileParams struct {
	Gender  enum.Gender
	Mobile  string
	Email   string
	Address string
}

func FindUser(ctx context.Context, uid uint64) (*models.User, error) {
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, preloadAvatar(ctx, user)
}

func UpdateUserProfile(ctx context.Context, uid uint64, params *UserProfileParams) (*models.User, error) {
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	user.Gender = params.Gender
	user.Mobile = params.Mobile
	user.Email = params.Email
	user.Address = params.Address
	if err = models.UpdateUserProfile(ctx, user); err != nil {
		return nil, err
	}
	return user, preloadAvatar(ctx, user)
}

func UpdateUserAvatar(ctx context.Context, uid, attachmentID uint64) (*models.User, error) {
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	user.AvatarID = attachmentID
	if err = models.UpdateUserAvatar(ctx, user); err != nil {
		return nil, err
	}
	return user, preloadAvatar(ctx, user)
}

func UpdateUsername(ctx context.Context, uid uint64, username string) (*models.User, error) {
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	if user.Username != "" {
		return nil, codes.ErrUsernameExisted
	}
	user.Username = username
	if err = models.UpdateUsername(ctx, user); err != nil {
		return nil, err
	}
	return user, preloadAvatar(ctx, user)
}

func preloadAvatar(ctx context.Context, users ...*models.User) error {
	avatarIDs := make([]uint64, len(users))
	for i, user := range users {
		avatarIDs[i] = user.AvatarID
	}
	attachmentMap, err := models.FindAttachmentMap(ctx, avatarIDs)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.Avatar = attachmentMap[user.AvatarID]
	}
	return nil
}
