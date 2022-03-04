package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mises-id/sns-socialsvc/app/factory"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	"github.com/mises-id/sns-socialsvc/app/models/meta"
	airdropSVC "github.com/mises-id/sns-socialsvc/app/services/airdrop"
	blacklistSVC "github.com/mises-id/sns-socialsvc/app/services/blacklist"
	commentSVC "github.com/mises-id/sns-socialsvc/app/services/comment"
	friendshipSVC "github.com/mises-id/sns-socialsvc/app/services/follow"
	messageSVC "github.com/mises-id/sns-socialsvc/app/services/message"
	sessionSVC "github.com/mises-id/sns-socialsvc/app/services/session"
	statusSVC "github.com/mises-id/sns-socialsvc/app/services/status"
	userSVC "github.com/mises-id/sns-socialsvc/app/services/user"
	twitterSVC "github.com/mises-id/sns-socialsvc/app/services/user_twitter"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	"github.com/mises-id/sns-socialsvc/lib/utils"
	pb "github.com/mises-id/sns-socialsvc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type requestWithCurrentUID interface {
	GetCurrentUid() uint64
}

func contextWithCurrentUID(parent context.Context, in requestWithCurrentUID) context.Context {
	if in.GetCurrentUid() == 0 {
		return parent
	}
	return context.WithValue(parent, utils.CurrentUIDKey{}, in.GetCurrentUid())
}

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.SocialServer {
	return socialService{}
}

type socialService struct{}

func (s socialService) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	var resp pb.SignInResponse
	jwt, created, err := sessionSVC.SignIn(ctx, in.Auth)
	if err != nil {
		return nil, err
	}
	resp.Jwt = jwt
	resp.IsCreated = created
	return &resp, nil
}

func (s socialService) FindUser(ctx context.Context, in *pb.FindUserRequest) (*pb.FindUserResponse, error) {
	var resp pb.FindUserResponse
	user, err := userSVC.FindUser(contextWithCurrentUID(ctx, in), in.Uid)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.User = factory.NewUserInfo(user)
	resp.IsFollowed = user.IsFollowed
	return &resp, nil
}

func (s socialService) UpdateUserProfile(ctx context.Context, in *pb.UpdateUserProfileRequest) (*pb.UpdateUserResponse, error) {
	var resp pb.UpdateUserResponse
	gender, err := enum.GenderFromString(in.Gender)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	user, err := userSVC.UpdateUserProfile(ctx, in.Uid, &userSVC.UserProfileParams{
		Gender:  gender,
		Mobile:  in.Mobile,
		Email:   in.Email,
		Address: in.Address,
	})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.User = factory.NewUserInfo(user)

	return &resp, nil
}

func (s socialService) UpdateUserAvatar(ctx context.Context, in *pb.UpdateUserAvatarRequest) (*pb.UpdateUserResponse, error) {
	var resp pb.UpdateUserResponse
	user, err := userSVC.UpdateUserAvatar(ctx, in.Uid, in.AttachmentPath)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.User = factory.NewUserInfo(user)
	return &resp, nil
}

func (s socialService) UpdateUserName(ctx context.Context, in *pb.UpdateUserNameRequest) (*pb.UpdateUserResponse, error) {
	var resp pb.UpdateUserResponse
	user, err := userSVC.UpdateUsername(ctx, in.Uid, in.Username)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.User = factory.NewUserInfo(user)
	return &resp, nil
}

func (s socialService) GetStatus(ctx context.Context, in *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	var resp pb.GetStatusResponse
	statusID, err := primitive.ObjectIDFromHex(in.Statusid)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}

	status, err := statusSVC.GetStatus(ctx, in.CurrentUid, statusID)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Status = factory.NewStatusInfo(status)
	return &resp, nil
}

func (s socialService) ListStatus(ctx context.Context, in *pb.ListStatusRequest) (*pb.ListStatusResponse, error) {
	var resp pb.ListStatusResponse
	fromTypes := []enum.FromType{}
	for _, from := range in.FromTypes {
		fromType, err := enum.FromTypeFromString(from)
		if err != nil {
			return nil, codes.ErrInvalidArgument
		}
		fromTypes = append(fromTypes, fromType)
	}
	statuses, page, err := statusSVC.ListStatus(ctx, &statusSVC.ListStatusParams{
		PageQuickParams: &pagination.PageQuickParams{
			Limit:  int64(in.Paginator.Limit),
			NextID: in.Paginator.NextId,
		},
		CurrentUID: in.GetCurrentUid(),
		UID:        in.TargetUid,
		FromTypes:  fromTypes,
	})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Statuses = factory.NewStatusInfoSlice(statuses)
	quickpage := page.BuildJSONResult().(*pagination.QuickPagination)
	resp.Paginator = &pb.PageQuick{
		Limit:  uint64(quickpage.Limit),
		NextId: quickpage.NextID,
	}
	return &resp, nil
}

func (s socialService) CreateStatus(ctx context.Context, in *pb.CreateStatusRequest) (*pb.CreateStatusResponse, error) {
	var resp pb.CreateStatusResponse

	param := &statusSVC.CreateStatusParams{
		StatusType:   in.StatusType,
		Content:      in.Content,
		IsPrivate:    in.IsPrivate,
		ShowDuration: int64(in.ShowDuration),
	}
	fmt.Println(in)
	fromType, err := enum.FromTypeFromString(in.FromType)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	param.FromType = fromType
	if len(in.ParentId) > 0 {
		parentID, err := primitive.ObjectIDFromHex(in.ParentId)
		if err != nil {
			return nil, codes.ErrInvalidArgument
		}
		param.ParentID = parentID
	}
	statusType, err := enum.StatusTypeFromString(in.StatusType)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	var data meta.MetaData
	switch statusType {
	default:
		data.TextMeta = &meta.TextMeta{}
	case enum.TextStatus:
		data.TextMeta = &meta.TextMeta{}
		_ = json.Unmarshal([]byte(in.Meta), data.TextMeta)
	case enum.LinkStatus:
		data.LinkMeta = &meta.LinkMeta{}
		_ = json.Unmarshal([]byte(in.Meta), data.LinkMeta)
	case enum.ImageStatus:
		data.ImageMeta = &meta.ImageMeta{Images: in.Images}
	}
	param.Meta = data
	status, err := statusSVC.CreateStatus(ctx, in.CurrentUid, param)

	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Status = factory.NewStatusInfo(status)
	return &resp, nil
}

func (s socialService) UpdateStatus(ctx context.Context, in *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {
	var resp pb.UpdateStatusResponse
	statusID, err := primitive.ObjectIDFromHex(in.StatusId)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	param := &statusSVC.UpdateStatusParams{
		ID:           statusID,
		IsPrivate:    in.IsPrivate,
		ShowDuration: int64(in.ShowDuration),
	}
	status, err := statusSVC.UpdateStatus(ctx, in.CurrentUid, param)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Status = factory.NewStatusInfo(status)
	return &resp, nil
}

func (s socialService) DeleteStatus(ctx context.Context, in *pb.DeleteStatusRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	statusID, err := primitive.ObjectIDFromHex(in.Statusid)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	if err := statusSVC.DeleteStatus(ctx, in.CurrentUid, statusID); err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}

func (s socialService) UnLikeStatus(ctx context.Context, in *pb.UnLikeStatusRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	statusID, err := primitive.ObjectIDFromHex(in.Statusid)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	if err := statusSVC.UnlikeStatus(ctx, in.CurrentUid, statusID); err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}

func (s socialService) LikeStatus(ctx context.Context, in *pb.LikeStatusRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	statusID, err := primitive.ObjectIDFromHex(in.Statusid)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	if _, err := statusSVC.LikeStatus(ctx, in.CurrentUid, statusID); err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}

func (s socialService) ListUserTimeline(ctx context.Context, in *pb.ListStatusRequest) (*pb.ListStatusResponse, error) {
	var resp pb.ListStatusResponse
	statuses, page, err := statusSVC.UserTimeline(ctx, in.CurrentUid, &pagination.PageQuickParams{
		Limit:  int64(in.Paginator.Limit),
		NextID: in.Paginator.NextId,
	})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Statuses = factory.NewStatusInfoSlice(statuses)
	quickpage := page.BuildJSONResult().(*pagination.QuickPagination)
	resp.Paginator = &pb.PageQuick{
		Limit:  uint64(quickpage.Limit),
		NextId: quickpage.NextID,
	}

	return &resp, nil
}

func (s socialService) ListRelationship(ctx context.Context, in *pb.ListRelationshipRequest) (*pb.ListRelationshipResponse, error) {
	var resp pb.ListRelationshipResponse
	relationType, err := enum.RelationTypeFromString(in.RelationType)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	relations, page, err := friendshipSVC.ListFriendship(contextWithCurrentUID(ctx, in), in.Uid, relationType, &pagination.QuickPagination{
		Limit:  int64(in.Paginator.Limit),
		NextID: in.Paginator.NextId,
	})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Relations = factory.NewRelationInfoSlice(relationType, relations)
	quickpage := page.BuildJSONResult().(*pagination.QuickPagination)
	resp.Paginator = &pb.PageQuick{
		Limit:  uint64(quickpage.Limit),
		NextId: quickpage.NextID,
	}
	return &resp, nil
}

func (s socialService) ListRecommended(ctx context.Context, in *pb.ListStatusRequest) (*pb.ListStatusResponse, error) {
	var resp pb.ListStatusResponse
	statuses, page, err := statusSVC.RecommendStatus(ctx, in.CurrentUid, &pagination.PageQuickParams{
		Limit:  int64(in.Paginator.Limit),
		NextID: in.Paginator.NextId,
	})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Statuses = factory.NewStatusInfoSlice(statuses)
	quickpage := page.BuildJSONResult().(*pagination.QuickPagination)
	resp.Paginator = &pb.PageQuick{
		Limit:  uint64(quickpage.Limit),
		NextId: quickpage.NextID,
	}

	return &resp, nil
}

func (s socialService) UnFollow(ctx context.Context, in *pb.UnFollowRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	if err := friendshipSVC.Unfollow(ctx, in.CurrentUid, in.TargetUid); err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}

func (s socialService) Follow(ctx context.Context, in *pb.FollowRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	if _, err := friendshipSVC.Follow(ctx, in.CurrentUid, in.TargetUid); err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}

func (s socialService) ListMessage(ctx context.Context, in *pb.ListMessageRequest) (*pb.ListMessageResponse, error) {
	var resp pb.ListMessageResponse
	messages, page, err := messageSVC.ListMessage(ctx, &messageSVC.ListMessageParams{
		ListMessageParams: models.ListMessageParams{
			UID: in.GetCurrentUid(),
			PageParams: &pagination.PageQuickParams{
				Limit:  int64(in.Paginator.Limit),
				NextID: in.Paginator.NextId,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	resp.Code = 0
	resp.Messages = factory.NewMessageSlice(messages)
	quickpage := page.BuildJSONResult().(*pagination.QuickPagination)
	resp.Paginator = &pb.PageQuick{
		Limit:  uint64(quickpage.Limit),
		NextId: quickpage.NextID,
	}
	return &resp, nil
}

func (s socialService) ReadMessage(ctx context.Context, in *pb.ReadMessageRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	var latestID primitive.ObjectID
	var messageIDs []primitive.ObjectID
	var err error
	if in.GetIds() != nil {
		messageIDs = make([]primitive.ObjectID, len(in.GetIds()))
		for i, id := range in.GetIds() {
			messageIDs[i], err = primitive.ObjectIDFromHex(id)
			if err != nil {
				return nil, err
			}
		}
	}
	if in.GetLatestID() != "" {
		latestID, err = primitive.ObjectIDFromHex(in.GetLatestID())
		if err != nil {
			return nil, err
		}
	}
	err = messageSVC.ReadMessages(ctx, &messageSVC.ReadMessageParams{
		ReadMessageParams: &models.ReadMessageParams{
			UID:        in.GetCurrentUid(),
			MessageIDs: messageIDs,
			LatestID:   latestID,
		},
	})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}

func (s socialService) ListComment(ctx context.Context, in *pb.ListCommentRequest) (*pb.ListCommentResponse, error) {
	var resp pb.ListCommentResponse
	statusID, err := primitive.ObjectIDFromHex(in.GetStatusId())
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	var groupID primitive.ObjectID
	if in.GetTopicId() != "" {
		groupID, err = primitive.ObjectIDFromHex(in.GetTopicId())
		if err != nil {
			return nil, codes.ErrInvalidArgument
		}
	}
	comments, page, err := commentSVC.ListComment(contextWithCurrentUID(ctx, in), &commentSVC.ListCommentParams{
		ListCommentParams: models.ListCommentParams{
			StatusID: statusID,
			GroupID:  groupID,
			PageParams: &pagination.PageQuickParams{
				Limit:  int64(in.Paginator.Limit),
				NextID: in.Paginator.NextId,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var total uint64
	status, err := statusSVC.GetStatusData(ctx, in.CurrentUid, statusID)
	if err == nil {
		total = status.CommentsCount
	}
	resp.Code = 0
	resp.Comments = factory.NewCommentSlice(comments)
	quickpage := page.BuildJSONResult().(*pagination.QuickPagination)
	resp.Paginator = &pb.PageQuick{
		Limit:  uint64(quickpage.Limit),
		NextId: quickpage.NextID,
		Total:  total,
	}
	return &resp, nil
}

func (s socialService) LatestFollowing(ctx context.Context, in *pb.LatestFollowingRequest) (*pb.LatestFollowingResponse, error) {
	var resp pb.LatestFollowingResponse
	followings, err := friendshipSVC.LatestFollowing(ctx, in.CurrentUid)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Followings = factory.NewFollowingSlice(followings)

	return &resp, nil
}

func (s socialService) CreateComment(ctx context.Context, in *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	var resp pb.CreateCommentResponse
	statusID, err := primitive.ObjectIDFromHex(in.GetStatusId())
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	var parentID primitive.ObjectID
	if in.GetParentId() != "" {
		parentID, err = primitive.ObjectIDFromHex(in.GetParentId())
		if err != nil {
			return nil, codes.ErrInvalidArgument
		}
	}
	comment, err := commentSVC.CreateComment(ctx, &commentSVC.CreateCommentParams{
		CreateCommentParams: &models.CreateCommentParams{
			StatusID: statusID,
			ParentID: parentID,
			UID:      in.GetCurrentUid(),
			Content:  in.GetContent(),
		},
	})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Comment = factory.NewComment(comment)
	return &resp, nil
}

func (s socialService) LikeComment(ctx context.Context, in *pb.LikeCommentRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	commentID, err := primitive.ObjectIDFromHex(in.GetCommentId())
	if err != nil {
		return nil, err
	}
	if _, err := commentSVC.LikeComment(ctx, in.CurrentUid, commentID); err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}

func (s socialService) UnlikeComment(ctx context.Context, in *pb.UnlikeCommentRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	commentID, err := primitive.ObjectIDFromHex(in.GetCommentId())
	if err != nil {
		return nil, err
	}
	if err := commentSVC.UnlikeComment(ctx, in.CurrentUid, commentID); err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}

func (s socialService) ListLikeStatus(ctx context.Context, in *pb.ListLikeRequest) (*pb.ListLikeResponse, error) {
	var resp pb.ListLikeResponse
	likes, page, err := statusSVC.ListLikeStatus(contextWithCurrentUID(ctx, in), &statusSVC.ListLikeStatusParams{
		UID: in.GetUid(),
		PageParams: &pagination.PageQuickParams{
			Limit:  int64(in.Paginator.Limit),
			NextID: in.Paginator.NextId,
		},
	})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Statuses = factory.NewStatusLikeSlice(likes)
	quickpage := page.BuildJSONResult().(*pagination.QuickPagination)
	resp.Paginator = &pb.PageQuick{
		Limit:  uint64(quickpage.Limit),
		NextId: quickpage.NextID,
	}
	return &resp, nil
}

func (s socialService) DeleteBlacklist(ctx context.Context, in *pb.DeleteBlacklistRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	err := blacklistSVC.DeleteBlacklist(ctx, in.GetUid(), in.GetTargetUid())
	resp.Code = 0
	return &resp, err
}

func (s socialService) ListBlacklist(ctx context.Context, in *pb.ListBlacklistRequest) (*pb.ListBlacklistResponse, error) {
	var resp pb.ListBlacklistResponse
	blacklists, page, err := blacklistSVC.ListBlacklist(ctx, &blacklistSVC.ListBlacklistParams{
		UID: in.GetUid(),
		PageParams: &pagination.PageQuickParams{
			Limit:  int64(in.Paginator.Limit),
			NextID: in.Paginator.NextId,
		},
	})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Blacklists = factory.NewBlacklistSlice(blacklists)
	quickpage := page.BuildJSONResult().(*pagination.QuickPagination)
	resp.Paginator = &pb.PageQuick{
		Limit:  uint64(quickpage.Limit),
		NextId: quickpage.NextID,
	}

	return &resp, nil
}

func (s socialService) CreateBlacklist(ctx context.Context, in *pb.CreateBlacklistRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	_, err := blacklistSVC.CreateBlacklist(ctx, in.GetUid(), in.GetTargetUid())
	resp.Code = 0
	return &resp, err
}

func (s socialService) GetMessageSummary(ctx context.Context, in *pb.GetMessageSummaryRequest) (*pb.MessageSummaryResponse, error) {
	var resp pb.MessageSummaryResponse
	summary, err := messageSVC.GetMessageSummary(ctx, in.GetCurrentUid())
	if err != nil {
		return nil, err
	}
	resp.Summary = &pb.MessageSummary{
		LatestMessage:      factory.NewMessage(summary.LatestMessage),
		Total:              summary.Total,
		NotificationsCount: summary.NotificationsCount,
		UsersCount:         summary.UsersCount,
	}
	resp.Code = 0
	return &resp, nil
}

func (s socialService) NewRecommendStatus(ctx context.Context, in *pb.NewRecommendStatusRequest) (*pb.NewRecommendStatusResponse, error) {
	var resp pb.NewRecommendStatusResponse
	svcin := &statusSVC.NewRecommendInput{
		LastRecommendTime: int64(in.LastRecommendTime),
		LastCommonTime:    int64(in.LastCommonTime),
	}
	svcout, err := statusSVC.NewRecommendStatus(ctx, in.CurrentUid, svcin)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Statuses = factory.NewStatusInfoSlice(svcout.Data)
	resp.Num = uint64(len(svcout.Data))
	resp.Next = &pb.NewRecommendNext{
		LastRecommendTime: (svcout.Next.LastRecommendTime),
		LastCommonTime:    (svcout.Next.LastCommonTime),
	}

	return &resp, nil
}

func (s socialService) DeleteComment(ctx context.Context, in *pb.DeleteCommentRequest) (*pb.SimpleResponse, error) {
	var resp pb.SimpleResponse
	commentID, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	if err := commentSVC.DeleteComment(ctx, in.CurrentUid, commentID); err != nil {
		return nil, err
	}
	resp.Code = 0
	return &resp, nil
}

func (s socialService) GetComment(ctx context.Context, in *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	var resp pb.GetCommentResponse
	commentID, err := primitive.ObjectIDFromHex(in.CommentId)
	if err != nil {
		return nil, codes.ErrInvalidArgument
	}
	comment, err := commentSVC.GetComment(ctx, in.CurrentUid, commentID)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Comment = factory.NewComment(comment)
	return &resp, nil
}

func (s socialService) NewListStatus(ctx context.Context, in *pb.NewListStatusRequest) (*pb.NewListStatusResponse, error) {
	var resp pb.NewListStatusResponse
	fromTypes := []enum.FromType{}
	for _, from := range in.FromTypes {
		fromType, err := enum.FromTypeFromString(from)
		if err != nil {
			return nil, codes.ErrInvalidArgument
		}
		fromTypes = append(fromTypes, fromType)
	}
	ids := []primitive.ObjectID{}
	for _, id := range in.Ids {
		Id, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			ids = append(ids, Id)
		}
	}
	svcin := &statusSVC.NewListStatusInput{
		CurrentUID: in.CurrentUid,
		ListNum:    int64(in.ListNum),
		FromTypes:  fromTypes,
		IDs:        ids,
	}
	svcout, err := statusSVC.NewListStatus(ctx, svcin)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Statuses = factory.NewStatusInfoSlice(svcout)
	return &resp, nil
}

func (s socialService) ShareTweetUrl(ctx context.Context, in *pb.ShareTweetUrlRequest) (*pb.ShareTweetUrlResponse, error) {
	var resp pb.ShareTweetUrlResponse
	url, err := twitterSVC.GetShareTweetUrl(ctx, in.CurrentUid)
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Url = url
	return &resp, nil
}

func (s socialService) UserTwitterAuth(ctx context.Context, in *pb.UserTwitterAuthRequest) (*pb.UserTwitterAuthResponse, error) {
	var resp pb.UserTwitterAuthResponse
	twitterSVC.UserTwitterAuth()
	return &resp, nil
}

func (s socialService) UserTwitterAirdrop(ctx context.Context, in *pb.UserTwitterAirdropRequest) (*pb.UserTwitterAirdropResponse, error) {
	var resp pb.UserTwitterAirdropResponse
	airdropSVC.TwitterAirdrop(ctx)
	return &resp, nil
}
