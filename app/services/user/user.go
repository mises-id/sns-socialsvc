package user

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserProfileParams struct {
	Gender  enum.Gender
	Mobile  string
	Email   string
	Address string
	Intro   string
}
type UserConfig struct {
	NftState bool
}

func FindUser(ctx context.Context, uid uint64) (*models.User, error) {
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	if err = models.PreloadUserData(ctx, user); err != nil {
		return nil, err
	}
	models.UserMergeUserExt(ctx, user)
	user.AirdropStatus = models.GetAirdropStatus(ctx)
	user.NewFansCount, err = models.NewFansCount(ctx, uid)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func FindMisesUser(ctx context.Context, misesid string) (*models.User, error) {
	user, err := models.FindUserByMisesid(ctx, misesid)
	if err != nil {
		return nil, err
	}
	if err = models.PreloadUserData(ctx, user); err != nil {
		return nil, err
	}
	models.UserMergeUserExt(ctx, user)
	user.AirdropStatus = models.GetAirdropStatus(ctx)
	user.NewFansCount, err = models.NewFansCount(ctx, user.UID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserConfig(ctx context.Context, currentUID uint64, in *UserConfig) (*UserConfig, error) {
	return GetUserConfig(ctx, currentUID, currentUID)
}

func GetUserConfig(ctx context.Context, currentUID uint64, uid uint64) (*UserConfig, error) {
	user_config := &UserConfig{
		NftState: true,
	}
	return user_config, nil
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
	user.Intro = params.Intro
	if err = models.UpdateUserProfile(ctx, user); err != nil {
		return nil, err
	}
	return user, preloadAvatar(ctx, user)
}

func UpdateUserAvatar(ctx context.Context, uid uint64, attachmentPath string, nft_asset_id primitive.ObjectID) (*models.User, error) {
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	if !nft_asset_id.IsZero() {
		nft_asset, err := models.FindNftAssetByID(ctx, nft_asset_id)
		if err != nil {
			return nil, err
		}
		if nft_asset.UID != uid {
			return nil, codes.ErrForbidden
		}
		nft_avatar := &models.NftAvatar{
			NftAssetID:        nft_asset.ID,
			ImageURL:          nft_asset.ImageURL,
			ImagePreviewUrl:   nft_asset.ImagePreviewUrl,
			ImageThumbnailUrl: nft_asset.ImageThumbnailUrl,
		}
		user.NftAvatar = nft_avatar
		if err = models.UpdateUserNftAvatar(ctx, user); err != nil {
			return nil, err
		}
	} else {
		user.AvatarPath = attachmentPath
		user.NftAvatar = nil
		if err = models.UpdateUserAvatar(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, preloadAvatar(ctx, user)
}

func UpdateUsername(ctx context.Context, uid uint64, username string) (*models.User, error) {
	user, err := models.FindUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	user.Username = username
	if err = models.UpdateUsername(ctx, user); err != nil {
		return nil, err
	}
	return user, preloadAvatar(ctx, user)
}

func preloadAvatar(ctx context.Context, users ...*models.User) error {
	return models.PreloadUserAvatar(ctx, users...)
	/* paths := make([]string, 0)
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
	return nil */
}
