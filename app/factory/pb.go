package factory

import (
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/message"
	"github.com/mises-id/sns-socialsvc/app/models/meta"
	"github.com/mises-id/sns-socialsvc/lib/utils"
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
		Intro:           user.Intro,
		IsFollowed:      user.IsFollowed,
		IsAirdropped:    user.IsAirdropped,
		AirdropStatus:   user.AirdropStatus,
		IsBlocked:       user.IsBlocked,
		FollowingsCount: user.FollowingCount,
		FansCount:       user.FansCount,
		LikedCount:      user.LikedCount,
		NewFansCount:    user.NewFansCount,
		IsLogined:       user.IsLogined,
		HelpMisesid:     user.Misesid,
	}
	if user.NftAvatar != nil {
		userinfo.AvatarUrl = &pb.UserAvatar{
			Small:      user.NftAvatar.ImageThumbnailUrl,
			Medium:     user.NftAvatar.ImagePreviewUrl,
			Large:      user.NftAvatar.ImageURL,
			NftAssetId: user.NftAvatar.NftAssetID.Hex(),
		}
	} else {
		userinfo.AvatarUrl = &pb.UserAvatar{
			Small:      user.AvatarUrl,
			Medium:     user.AvatarUrl,
			Large:      user.AvatarUrl,
			NftAssetId: "",
		}
		if user.Avatar != nil {
			userinfo.AvatarUrl.Small = user.Avatar.Small
			userinfo.AvatarUrl.Medium = user.Avatar.Medium
		}
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
func NewAirdrop(in *models.Airdrop) *pb.Airdrop {
	if in == nil {
		return nil
	}
	out := &pb.Airdrop{
		Coin:      float32(utils.UMisesToMises(uint64(in.Coin))),
		Status:    in.Status.String(),
		FinishAt:  uint64(in.FinishAt.Unix()),
		CreatedAt: uint64(in.CreatedAt.Unix()),
	}
	return out
}

func NewUserTwitterAuth(in *models.UserTwitterAuth) *pb.UserTwitterAuth {
	if in == nil {
		return nil
	}
	out := &pb.UserTwitterAuth{
		TwitterUserId:    in.TwitterUserId,
		Misesid:          utils.RemoveMisesidProfix(in.Misesid),
		Name:             in.TwitterUser.Name,
		Username:         in.TwitterUser.UserName,
		FollowersCount:   in.TwitterUser.FollowersCount,
		TweetCount:       in.TwitterUser.TweetCount,
		TwitterCreatedAt: uint64(in.TwitterUser.CreatedAt.Unix()),
		Amount:           float32(utils.UMisesToMises(uint64(in.Amount))),
		CreatedAt:        uint64(in.CreatedAt.Unix()),
	}
	return out
}

func NewImageMetaInfo(meta *meta.ImageMeta) *pb.ImageMetaInfo {
	if meta == nil {
		return nil
	}
	info := &pb.ImageMetaInfo{
		Images:      meta.ImageURLs,
		ThumbImages: meta.ThumbImageURLs,
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
		ParentStatusIsBlacked: status.ParentStatusIsBlacked,
		CreatedAt:             uint64(status.CreatedAt.Unix()),
		IsPublic:              status.IsPublic,
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
func NewNftAssetSlice(assets []*models.NftAsset) []*pb.NftAsset {
	result := make([]*pb.NftAsset, len(assets))
	for i, asset := range assets {
		result[i] = NewNftAsset(asset)
	}
	return result
}
func NewNftEventSlice(events []*models.NftEvent) []*pb.NftEvent {
	result := make([]*pb.NftEvent, len(events))
	for i, event := range events {
		result[i] = NewNftEvent(event)
	}
	return result
}
func NewLikeSlice(likes []*models.Like) []*pb.Like {
	result := make([]*pb.Like, len(likes))
	for i, like := range likes {
		result[i] = NewLike(like)
	}
	return result
}
func NewLike(like *models.Like) *pb.Like {
	if like == nil {
		return nil
	}
	resp := &pb.Like{
		Id: docID(like.ID),
	}
	resp.User = NewUserInfo(like.User)
	return resp
}

func NewAssetContract(in *models.AssetContract) *pb.AssetContract {
	if in == nil {
		return nil
	}
	result := &pb.AssetContract{
		Address: in.Address,
	}
	return result
}
func NewCollection(in *models.NftCollection) *pb.NftCollection {
	if in == nil {
		return nil
	}
	result := &pb.NftCollection{
		Name:  in.Name,
		Slug:  in.Slug,
		Stats: NewStats(in.Stats),
	}
	if in.PaymentToken != nil {
		result.PaymentToken = NewPaymentTokenSlice(in.PaymentToken)
	}
	return result
}

func NewPaymentTokenSlice(in []*models.PaymentToken) []*pb.PaymentToken {
	result := make([]*pb.PaymentToken, len(in))
	for i, v := range in {
		result[i] = NewPaymentToken(v)
	}
	return result
}

func NewPaymentToken(in *models.PaymentToken) *pb.PaymentToken {
	if in == nil {
		return nil
	}
	result := &pb.PaymentToken{
		Id:       uint64(in.ID),
		Symbol:   in.Symbol,
		Address:  in.Address,
		Name:     in.Name,
		EthPrice: in.ETHPrice,
		UsdPrice: in.USDPrice,
		Decimals: in.Decimals,
	}
	return result
}

func NewStats(in *models.Stats) *pb.Stats {
	if in == nil {
		return nil
	}
	result := &pb.Stats{
		FloorPrice: float32(in.FloorPrice),
	}
	return result
}

func NewNftAsset(asset *models.NftAsset) *pb.NftAsset {
	if asset == nil {
		return nil
	}
	result := &pb.NftAsset{
		Id:                docID(asset.ID),
		ImageUrl:          asset.ImageURL,
		ImagePreviewUrl:   asset.ImagePreviewUrl,
		ImageThumbnailUrl: asset.ImageThumbnailUrl,
		TokenId:           asset.TokenId,
		PermaLink:         asset.PermaLink,
		LikesCount:        asset.LikesCount,
		CommentsCount:     asset.CommentsCount,
		Name:              asset.Name,
		User:              NewUserInfo(asset.User),
		IsLiked:           asset.IsLiked,
	}
	result.AssetContract = NewAssetContract(asset.AssetContract)
	result.Collection = NewCollection(asset.Collection)
	return result
}

func NewNftEvent(event *models.NftEvent) *pb.NftEvent {
	if event == nil {
		return nil
	}
	result := &pb.NftEvent{
		Id:           docID(event.ID),
		EventType:    event.EventType,
		FromAccount:  NewNftAccount(event.FromAccount),
		ToAccount:    NewNftAccount(event.ToAccount),
		CreatedDate:  event.CreatedDate,
		PaymentToken: NewPaymentToken(event.PaymentToken),
	}
	return result
}

func NewNftAccount(account *models.Account) *pb.NftAccount {
	if account == nil {
		return nil
	}
	result := &pb.NftAccount{
		Address:       account.Address,
		ProfileImgUrl: account.ProfileImgUrl,
		MisesUser:     NewUserInfo(account.MisesUser),
	}
	return result
}

func NewRelationInfoSlice(relationType enum.RelationType, follows []*models.Follow) []*pb.RelationInfo {
	result := make([]*pb.RelationInfo, len(follows))
	for i, follow := range follows {
		user := follow.ToUser
		currentRelationType := enum.Fan
		if relationType == enum.Fan {
			user = follow.FromUser
			//currentRelationType = enum.Fan
		}
		/* if follow.IsFriend {
			currentRelationType = enum.Friend
		} */
		if user != nil {
			if user.IsFollowed {
				currentRelationType = enum.Following
			}
			if user.IsFriend {
				currentRelationType = enum.Friend
			}
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
func newNftCommentMeta(meta *message.NftAssetCommentMeta) *pb.NewNftCommentMeta {
	if meta == nil {
		return &pb.NewNftCommentMeta{}
	}
	return &pb.NewNftCommentMeta{
		Uid:            meta.UID,
		GroupId:        docID(meta.GroupID),
		CommentId:      docID(meta.CommentID),
		Content:        meta.Content,
		ParentContent:  meta.ParentContent,
		ParentUserName: meta.ParentUsername,
		NftAssetName:   meta.NftAssetName,
		NftAssetImage:  meta.NftAssetImage,
	}
}

func newLikeNftMeta(meta *message.LikeNftAssetMeta) *pb.NewLikeNftMeta {
	if meta == nil {
		return &pb.NewLikeNftMeta{}
	}
	return &pb.NewLikeNftMeta{
		Uid:           meta.UID,
		NftAssetId:    meta.NftAssetID.Hex(),
		NftAssetName:  meta.NftAssetName,
		NftAssetImage: meta.NftAssetImage,
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
func newLikeNftCommentMeta(meta *message.LikeNftAssetCommentMeta) *pb.NewLikeNftCommentMeta {
	if meta == nil {
		return &pb.NewLikeNftCommentMeta{}
	}
	return &pb.NewLikeNftCommentMeta{
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
		Id:               docID(message.ID),
		Uid:              message.UID,
		MessageType:      message.MessageType.String(),
		FromUser:         NewUserInfo(message.FromUser),
		State:            message.State(),
		Status:           NewStatusInfo(message.Status),
		NftAsset:         NewNftAsset(message.NftAsset),
		CreatedAt:        uint64(message.CreatedAt.Unix()),
		StatusIsDeleted:  message.StatusIsDeleted,
		CommentIsDeleted: message.CommentIsDeleted,
	}
	switch message.MessageType {
	case enum.NewComment:
		result.NewCommentMeta = newCommentMeta(message.CommentMeta)
	case enum.NewLikeStatus:
		result.NewLikeStatusMeta = newLikeStatusMeta(message.LikeStatusMeta)
	case enum.NewLikeComment:
		result.NewLikeCommentMeta = newLikeCommentMeta(message.LikeCommentMeta)
	case enum.NewNftComment:
		result.NewNftComment = newNftCommentMeta(message.NftCommentMeta)
	case enum.NewLikeNftComment:
		result.NewLikeNftCommentMeta = newLikeNftCommentMeta(message.LikeNftCommentMeta)
	case enum.NewLikeNft:
		result.NewLikeNft = newLikeNftMeta(message.LikeNftMeta)
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
func NewChannelUserListSlice(channel_users []*models.ChannelUser) []*pb.ChannelUserInfo {
	result := make([]*pb.ChannelUserInfo, len(channel_users))
	for i, channel_user := range channel_users {
		result[i] = NewChannelUser(channel_user)
	}
	return result
}

func NewChannelUser(channel_user *models.ChannelUser) *pb.ChannelUserInfo {

	return &pb.ChannelUserInfo{
		Id:             channel_user.ID.Hex(),
		ChannelId:      channel_user.ChannelID.Hex(),
		ValidState:     int32(channel_user.ValidState),
		Amount:         uint64(channel_user.Amount),
		TxId:           channel_user.TxID,
		User:           NewUserInfo(channel_user.User),
		AirdropState:   int32(channel_user.AirdropState),
		AirdropTime:    uint64(channel_user.AirdropTime.Unix()),
		CreatedAt:      uint64(channel_user.CreatedAt.Unix()),
		ChannelUid:     channel_user.ChannelUID,
		ChannelMisesid: channel_user.ChannelMisesid,
	}

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
