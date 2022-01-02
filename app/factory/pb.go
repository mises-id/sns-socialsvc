package factory

import (
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/message"
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

func newCommentMeta(meta *message.CommentMeta) *pb.Message_NewCommentMeta {
	return &pb.Message_NewCommentMeta{
		NewCommentMeta: &pb.NewCommentMeta{
			Uid:       meta.UID,
			GroupId:   meta.GroupID.Hex(),
			CommentId: meta.CommentID.Hex(),
			Content:   meta.Content,
		},
	}
}

func newLikeMeta(meta *message.LikeMeta) *pb.Message_NewLikeMeta {
	return &pb.Message_NewLikeMeta{
		NewLikeMeta: &pb.NewLikeMeta{
			Uid:        meta.UID,
			TargetId:   meta.TargetID.Hex(),
			TargetType: meta.TargetType.String(),
		},
	}
}

func newFansMeta(meta *message.FansMeta) *pb.Message_NewFansMeta {
	return &pb.Message_NewFansMeta{
		NewFansMeta: &pb.NewFansMeta{
			Uid: meta.UID,
		},
	}
}

func newForwardMeta(meta *message.ForwardMeta) *pb.Message_NewForwardMeta {
	return &pb.Message_NewForwardMeta{
		NewForwardMeta: &pb.NewForwardMeta{
			Uid:      meta.UID,
			StatusId: meta.StatusID.Hex(),
			Content:  meta.Content,
		},
	}
}
func NewMessage(message *models.Message) *pb.Message {
	if message == nil {
		return nil
	}
	result := &pb.Message{
		Id:          message.ID.Hex(),
		Uid:         message.UID,
		MessageType: message.MessageType.String(),
		State:       message.State(),
	}
	switch message.MessageType {
	case enum.NewComment:
		result.MetaData = newCommentMeta(message.CommentMeta)
	case enum.NewLike:
		result.MetaData = newLikeMeta(message.LikeMeta)
	case enum.NewFans:
		result.MetaData = newFansMeta(message.FansMeta)
	case enum.NewForward:
		result.MetaData = newForwardMeta(message.ForwardMeta)
	}
	return result
}

func NewMessageSlice(messages []*models.Message) []*pb.Message {
	result := make([]*pb.Message, len(messages))
	for i, message := range messages {
		result[i] = NewMessage(message)
	}
	return result
}

func NewComment(comment *models.Comment) *pb.Comment {
	if comment == nil {
		return nil
	}
	result := &pb.Comment{
		Id:         comment.ID.Hex(),
		Uid:        comment.UID,
		StatusId:   comment.StatusID.Hex(),
		ParentId:   comment.ParentID.Hex(),
		GroupId:    comment.GroupID.Hex(),
		OpponentId: comment.OpponentID,
		Content:    comment.Content,
	}
	return result
}

func NewCommentSlice(comments []*models.Comment) []*pb.Comment {
	result := make([]*pb.Comment, len(comments))
	for i, comment := range comments {
		result[i] = NewComment(comment)
	}
	return result
}
