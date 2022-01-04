package handlers

import (
	"context"

	"github.com/mises-id/sns-socialsvc/app/factory"
	"github.com/mises-id/sns-socialsvc/app/models"
	"github.com/mises-id/sns-socialsvc/app/models/enum"
	commentSVC "github.com/mises-id/sns-socialsvc/app/services/comment"
	friendshipSVC "github.com/mises-id/sns-socialsvc/app/services/follow"
	messageSVC "github.com/mises-id/sns-socialsvc/app/services/message"
	sessionSVC "github.com/mises-id/sns-socialsvc/app/services/session"
	statusSVC "github.com/mises-id/sns-socialsvc/app/services/status"
	userSVC "github.com/mises-id/sns-socialsvc/app/services/user"
	"github.com/mises-id/sns-socialsvc/lib/codes"
	"github.com/mises-id/sns-socialsvc/lib/pagination"
	pb "github.com/mises-id/sns-socialsvc/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.SocialServer {
	return socialService{}
}

type socialService struct{}

func (s socialService) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	var resp pb.SignInResponse
	jwt, err := sessionSVC.SignIn(ctx, in.Auth)
	if err != nil {
		return nil, err
	}
	resp.Jwt = jwt
	return &resp, nil
}

func (s socialService) FindUser(ctx context.Context, in *pb.FindUserRequest) (*pb.FindUserResponse, error) {
	var resp pb.FindUserResponse
	user, err := userSVC.FindUser(ctx, in.Uid)
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
		return nil, err
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
			return nil, err
		}
		fromTypes = append(fromTypes, fromType)
	}
	statuses, page, err := statusSVC.ListStatus(ctx, &statusSVC.ListStatusParams{
		PageQuickParams: &pagination.PageQuickParams{
			Limit:  int64(in.Paginator.Limit),
			NextID: in.Paginator.NextId,
		},
		CurrentUID: in.CurrentUid,
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
		StatusType: in.StatusType,
		Content:    in.Content,
	}
	fromType, err := enum.FromTypeFromString(in.FromType)
	if err != nil {
		return nil, err
	}
	param.FromType = fromType
	if len(in.ParentId) > 0 {
		parentID, err := primitive.ObjectIDFromHex(in.ParentId)
		if err != nil {
			return nil, err
		}
		param.ParentID = parentID
	}
	status, err := statusSVC.CreateStatus(ctx, in.CurrentUid, param)

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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}
	relations, page, err := friendshipSVC.ListFriendship(ctx, in.CurrentUid, relationType, &pagination.QuickPagination{
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
	return &resp, nil
}

func (s socialService) ListComment(ctx context.Context, in *pb.ListCommentRequest) (*pb.ListCommentResponse, error) {
	var resp pb.ListCommentResponse
	comments, page, err := commentSVC.ListComment(ctx, &commentSVC.ListCommentParams{
		ListCommentParams: models.ListCommentParams{},
	})
	if err != nil {
		return nil, err
	}
	resp.Code = 0
	resp.Comments = factory.NewCommentSlice(comments)
	quickpage := page.BuildJSONResult().(*pagination.QuickPagination)
	resp.Paginator = &pb.PageQuick{
		Limit:  uint64(quickpage.Limit),
		NextId: quickpage.NextID,
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
	return &resp, nil
}
