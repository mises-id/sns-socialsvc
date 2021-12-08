package socialsvc

import (
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/meta"
)

func NewUserInfo(user *models.User) *UserInfo {
	var avatarUrl = ""
	if user.Avatar != nil {
		avatarUrl = user.Avatar.FileUrl()
	}
	userinfo := UserInfo{
		Uid:      user.UID,
		Username: user.Username,
		Misesid:  user.Misesid,
		Gender:   user.Gender.String(),
		Mobile:   user.Mobile,
		Email:    user.Email,
		Address:  user.Address,
		Avatar:   avatarUrl,
	}
	return &userinfo
}

func NewLinkMetaInfo(meta *meta.LinkMeta) *LinkMetaInfo {
	if meta == nil {
		return nil
	}
	linkMetaInfo := LinkMetaInfo{
		Title:         meta.Title,
		Host:          meta.Host,
		Link:          meta.Link,
		AttachmentId:  meta.AttachmentID,
		AttachmentUrl: meta.AttachmentURL,
	}
	return &linkMetaInfo
}

func NewStatusInfo(status *models.Status) *StatusInfo {
	if status == nil {
		return nil
	}
	metaData, err := status.GetMetaData()
	if err != nil {
		return nil
	}

	statusinfo := StatusInfo{
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
		if metaData != nil {
			linkMeta := metaData.(*meta.LinkMeta)
			statusinfo.LinkMeta = NewLinkMetaInfo(linkMeta)
		}
	}
	return &statusinfo
}

func NewStatusInfoSlice(statuses []*models.Status) []*StatusInfo {
	result := make([]*StatusInfo, len(statuses))
	for i, status := range statuses {
		result[i] = NewStatusInfo(status)
	}
	return result
}

func NewRelationInfoSlice(relationType enum.RelationType, follows []*models.Follow) []*RelationInfo {
	result := make([]*RelationInfo, len(follows))
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
		result[i] = &RelationInfo{
			User:         NewUserInfo(user),
			RelationType: currentRelationType.String(),
			CreatedAt:    uint64(follow.CreatedAt.Unix()),
		}
	}
	return result
}
