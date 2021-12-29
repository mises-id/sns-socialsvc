package factory

import (
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/meta"
	pb "github.com/mises-id/sns-socialsvc/proto"
)

func NewUserInfo(user *models.User) *pb.UserInfo {
	userinfo := pb.UserInfo{
		Uid:      user.UID,
		Username: user.Username,
		Misesid:  user.Misesid,
		Gender:   user.Gender.String(),
		Mobile:   user.Mobile,
		Email:    user.Email,
		Address:  user.Address,
		Avatar:   user.AvatarUrl,
	}
	return &userinfo
}

func NewLinkMetaInfo(meta *meta.LinkMeta) *pb.LinkMetaInfo {
	if meta == nil {
		return nil
	}
	linkMetaInfo := pb.LinkMetaInfo{
		Title:     meta.Title,
		Host:      meta.Host,
		Link:      meta.Link,
		ImagePath: meta.ImagePath,
		ImageUrl:  meta.ImageURL,
	}
	return &linkMetaInfo
}

func NewImageMetaInfo(meta *meta.ImageMeta) *pb.ImageMetaInfo {
	if meta == nil {
		return nil
	}
	info := &pb.ImageMetaInfo{
		Images: meta.ImageURLs,
	}
	return info
}

func NewStatusInfo(status *models.Status) *pb.StatusInfo {
	if status == nil {
		return nil
	}

	statusinfo := pb.StatusInfo{
		Id:           status.ID.Hex(),
		User:         NewUserInfo(status.User),
		Content:      status.Content,
		FromType:     status.FromType.String(),
		StatusType:   status.StatusType.String(),
		Parent:       NewStatusInfo(status.ParentStatus),
		Origin:       NewStatusInfo(status.OriginStatus),
		CommentCount: status.CommentsCount,
		LikeCount:    status.LikesCount,
		ForwardCount: status.ForwardsCount,
		IsLiked:      status.IsLiked,
		CreatedAt:    uint64(status.CreatedAt.Unix()),
	}
	switch status.StatusType {
	case enum.LinkStatus:
		statusinfo.LinkMeta = NewLinkMetaInfo(status.LinkMeta)
	case enum.ImageStatus:
		statusinfo.ImageMeta = NewImageMetaInfo(status.ImageMeta)
	}
	return &statusinfo
}

func NewStatusInfoSlice(statuses []*models.Status) []*pb.StatusInfo {
	result := make([]*pb.StatusInfo, len(statuses))
	for i, status := range statuses {
		result[i] = NewStatusInfo(status)
	}
	return result
}

func NewRelationInfoSlice(relationType enum.RelationType, follows []*models.Follow) []*pb.RelationInfo {
	result := make([]*pb.RelationInfo, len(follows))
	for i, follow := range follows {
		user := follow.ToUser
		currentRelationType := enum.Following
		if relationType == enum.Fan {
			user = follow.FromUser
			currentRelationType = enum.Fan
		}
		if follow.IsFriend {
			currentRelationType = enum.Friend
		}
		result[i] = &pb.RelationInfo{
			User:         NewUserInfo(user),
			RelationType: currentRelationType.String(),
			CreatedAt:    uint64(follow.CreatedAt.Unix()),
		}
	}
	return result
}
