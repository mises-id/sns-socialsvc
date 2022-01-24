package factory

import (
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/message"
	"github.com/mises-id/sns-socialsvc/app/models/meta"
	pb "github.com/mises-id/sns-socialsvc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewUserInfo(user *models.User) *pb.UserInfo {
	if user == nil {
		return nil
	}
	userinfo := pb.UserInfo{
		Uid:             user.UID,
		Username:        user.Username,
		Misesid:         user.Misesid,
		Gender:          user.Gender.String(),
		Mobile:          user.Mobile,
		Email:           user.Email,
		Address:         user.Address,
		Avatar:          user.AvatarUrl,
		IsFollowed:      user.IsFollowed,
		IsBlocked:       user.IsBlocked,
		FollowingsCount: user.FollowingCount,
		FansCount:       user.FansCount,
		LikedCount:      user.LikedCount,
		NewFansCount:    user.NewFansCount,
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
		Id:                    docID(status.ID),
		User:                  NewUserInfo(status.User),
		Content:               status.Content,
		FromType:              status.FromType.String(),
		StatusType:            status.StatusType.String(),
		Parent:                NewStatusInfo(status.ParentStatus),
		Origin:                NewStatusInfo(status.OriginStatus),
		CommentCount:          status.CommentsCount,
		LikeCount:             status.LikesCount,
		ForwardCount:          status.ForwardsCount,
		IsLiked:               status.IsLiked,
		ParentStatusIsDeleted: status.ParentStatusIsDeleted,
		ParentStatusIsBlacked: status.ParentStatusIsBlocked,
		CreatedAt:             uint64(status.CreatedAt.Unix()),
		IsPublic:              status.HideTime == nil,
		HideTime:              status.GetHideTime(),
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

func newCommentMeta(meta *message.CommentMeta) *pb.NewCommentMeta {
	if meta == nil {
		return &pb.NewCommentMeta{}
	}
	return &pb.NewCommentMeta{
		Uid:                  meta.UID,
		GroupId:              docID(meta.GroupID),
		CommentId:            docID(meta.CommentID),
		Content:              meta.Content,
		ParentContent:        meta.ParentContent,
		ParentUserName:       meta.ParentUsername,
		StatusContentSummary: meta.StatusContentSummary,
		StatusImageUrl:       meta.StatusImageURL,
	}
}

func newLikeStatusMeta(meta *message.LikeStatusMeta) *pb.NewLikeStatusMeta {
	if meta == nil {
		return &pb.NewLikeStatusMeta{}
	}
	return &pb.NewLikeStatusMeta{
		Uid:            meta.UID,
		StatusId:       meta.StatusID.Hex(),
		StatusContent:  meta.StatusContent,
		StatusImageUrl: meta.StatusImageURL,
	}

}

func newLikeCommentMeta(meta *message.LikeCommentMeta) *pb.NewLikeCommentMeta {
	if meta == nil {
		return &pb.NewLikeCommentMeta{}
	}
	return &pb.NewLikeCommentMeta{
		Uid:             meta.UID,
		CommentId:       meta.CommentID.Hex(),
		CommentUsername: meta.CommentUsername,
		CommentContent:  meta.CommentContent,
	}
}

func newFansMeta(meta *message.FansMeta) *pb.NewFansMeta {
	if meta == nil {
		return &pb.NewFansMeta{}
	}
	return &pb.NewFansMeta{
		Uid:         meta.UID,
		FanUsername: meta.FanUsername,
	}
}

func newForwardMeta(meta *message.ForwardMeta) *pb.NewForwardMeta {
	if meta == nil {
		return &pb.NewForwardMeta{}
	}
	return &pb.NewForwardMeta{
		Uid:            meta.UID,
		StatusId:       meta.StatusID.Hex(),
		ForwardContent: meta.ForwardContent,
		ContentSummary: meta.ContentSummary,
		ImageUrl:       meta.ImageURL,
	}
}
func NewMessage(message *models.Message) *pb.Message {
	if message == nil {
		return nil
	}
	result := &pb.Message{
		Id:          docID(message.ID),
		Uid:         message.UID,
		MessageType: message.MessageType.String(),
		FromUser:    NewUserInfo(message.FromUser),
		State:       message.State(),
		Status:      NewStatusInfo(message.Status),
		CreatedAt:   uint64(message.CreatedAt.Unix()),
	}
	switch message.MessageType {
	case enum.NewComment:
		result.NewCommentMeta = newCommentMeta(message.CommentMeta)
	case enum.NewLikeStatus:
		result.NewLikeStatusMeta = newLikeStatusMeta(message.LikeStatusMeta)
	case enum.NewLikeComment:
		result.NewLikeCommentMeta = newLikeCommentMeta(message.LikeCommentMeta)
	case enum.NewFans:
		result.NewFansMeta = newFansMeta(message.FansMeta)
	case enum.NewForward:
		result.NewForwardMeta = newForwardMeta(message.ForwardMeta)
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
		Id:           docID(comment.ID),
		Uid:          comment.UID,
		StatusId:     docID(comment.StatusID),
		ParentId:     docID(comment.ParentID),
		GroupId:      docID(comment.GroupID),
		OpponentId:   comment.OpponentID,
		Content:      comment.Content,
		CommentCount: comment.CommentsCount,
		LikeCount:    comment.LikesCount,
		CreatedAt:    uint64(comment.CreatedAt.Unix()),
		IsLiked:      comment.IsLiked,
	}
	if comment.Comments != nil {
		result.Comments = NewCommentSlice(comment.Comments)
	}
	if comment.Opponent != nil {
		result.Opponent = NewUserInfo(comment.Opponent)
	}
	if comment.User != nil {
		result.User = NewUserInfo(comment.User)
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

func NewFollowingSlice(follows []*models.Follow) []*pb.Following {
	result := make([]*pb.Following, len(follows))
	for i, follow := range follows {
		result[i] = &pb.Following{
			User:   NewUserInfo(follow.ToUser),
			Unread: !follow.IsRead,
		}
	}
	return result
}

func NewBlacklistSlice(blacklists []*models.Blacklist) []*pb.Blacklist {
	result := make([]*pb.Blacklist, len(blacklists))
	for i, blacklist := range blacklists {
		result[i] = &pb.Blacklist{
			User:      NewUserInfo(blacklist.TargetUser),
			CreatedAt: uint64(blacklist.CreatedAt.Unix()),
		}
	}
	return result
}

func NewStatusLikeSlice(likes []*models.Like) []*pb.StatusLike {
	result := make([]*pb.StatusLike, len(likes))
	for i, like := range likes {
		result[i] = &pb.StatusLike{
			Status:    NewStatusInfo(like.Status),
			CreatedAt: uint64(like.CreatedAt.Unix()),
		}
	}
	return result
}

func docID(id primitive.ObjectID) string {
	if id.IsZero() {
		return ""
	}
	return id.Hex()
}
