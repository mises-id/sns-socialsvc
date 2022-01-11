package user

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/storage"
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
	if err = models.PreloadUserData(ctx, user); err != nil {
		return nil, err
	}
	user.NewFansCount, err = models.NewFansCount(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
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

func UpdateUserAvatar(ctx context.Context, uid uint64, attachmentPath string) (*models.User, error) {
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	user.AvatarPath = attachmentPath
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
	paths := make([]string, 0)
	for _, user := range users {
		if user.AvatarPath != "" {
			paths = append(paths, user.AvatarPath)
		}
	}
	avatars, err := storage.ImageClient.GetFileUrl(ctx, paths...)
	if err != nil {
		return err
	}
	for _, user := range users {
		user.AvatarUrl = avatars[user.AvatarPath]
	}
	return nil
}
