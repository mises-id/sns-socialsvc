// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 5f7d5bf015
// Version Date: 2021-11-26T09:27:01Z

package svc

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/mises-id/sns-socialsvc/proto"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC SocialServer.
func MakeGRPCServer(endpoints Endpoints, options ...grpctransport.ServerOption) pb.SocialServer {
	serverOptions := []grpctransport.ServerOption{
		grpctransport.ServerBefore(metadataToContext),
	}
	serverOptions = append(serverOptions, options...)
	return &grpcServer{
		// social

		signin: grpctransport.NewServer(
			endpoints.SignInEndpoint,
			DecodeGRPCSignInRequest,
			EncodeGRPCSignInResponse,
			serverOptions...,
		),
		finduser: grpctransport.NewServer(
			endpoints.FindUserEndpoint,
			DecodeGRPCFindUserRequest,
			EncodeGRPCFindUserResponse,
			serverOptions...,
		),
		updateuserprofile: grpctransport.NewServer(
			endpoints.UpdateUserProfileEndpoint,
			DecodeGRPCUpdateUserProfileRequest,
			EncodeGRPCUpdateUserProfileResponse,
			serverOptions...,
		),
		updateuseravatar: grpctransport.NewServer(
			endpoints.UpdateUserAvatarEndpoint,
			DecodeGRPCUpdateUserAvatarRequest,
			EncodeGRPCUpdateUserAvatarResponse,
			serverOptions...,
		),
		updateusername: grpctransport.NewServer(
			endpoints.UpdateUserNameEndpoint,
			DecodeGRPCUpdateUserNameRequest,
			EncodeGRPCUpdateUserNameResponse,
			serverOptions...,
		),
		createstatus: grpctransport.NewServer(
			endpoints.CreateStatusEndpoint,
			DecodeGRPCCreateStatusRequest,
			EncodeGRPCCreateStatusResponse,
			serverOptions...,
		),
		updatestatus: grpctransport.NewServer(
			endpoints.UpdateStatusEndpoint,
			DecodeGRPCUpdateStatusRequest,
			EncodeGRPCUpdateStatusResponse,
			serverOptions...,
		),
		deletestatus: grpctransport.NewServer(
			endpoints.DeleteStatusEndpoint,
			DecodeGRPCDeleteStatusRequest,
			EncodeGRPCDeleteStatusResponse,
			serverOptions...,
		),
		likestatus: grpctransport.NewServer(
			endpoints.LikeStatusEndpoint,
			DecodeGRPCLikeStatusRequest,
			EncodeGRPCLikeStatusResponse,
			serverOptions...,
		),
		unlikestatus: grpctransport.NewServer(
			endpoints.UnLikeStatusEndpoint,
			DecodeGRPCUnLikeStatusRequest,
			EncodeGRPCUnLikeStatusResponse,
			serverOptions...,
		),
		listlikestatus: grpctransport.NewServer(
			endpoints.ListLikeStatusEndpoint,
			DecodeGRPCListLikeStatusRequest,
			EncodeGRPCListLikeStatusResponse,
			serverOptions...,
		),
		getstatus: grpctransport.NewServer(
			endpoints.GetStatusEndpoint,
			DecodeGRPCGetStatusRequest,
			EncodeGRPCGetStatusResponse,
			serverOptions...,
		),
		liststatus: grpctransport.NewServer(
			endpoints.ListStatusEndpoint,
			DecodeGRPCListStatusRequest,
			EncodeGRPCListStatusResponse,
			serverOptions...,
		),
		listrecommended: grpctransport.NewServer(
			endpoints.ListRecommendedEndpoint,
			DecodeGRPCListRecommendedRequest,
			EncodeGRPCListRecommendedResponse,
			serverOptions...,
		),
		listusertimeline: grpctransport.NewServer(
			endpoints.ListUserTimelineEndpoint,
			DecodeGRPCListUserTimelineRequest,
			EncodeGRPCListUserTimelineResponse,
			serverOptions...,
		),
		latestfollowing: grpctransport.NewServer(
			endpoints.LatestFollowingEndpoint,
			DecodeGRPCLatestFollowingRequest,
			EncodeGRPCLatestFollowingResponse,
			serverOptions...,
		),
		listrelationship: grpctransport.NewServer(
			endpoints.ListRelationshipEndpoint,
			DecodeGRPCListRelationshipRequest,
			EncodeGRPCListRelationshipResponse,
			serverOptions...,
		),
		follow: grpctransport.NewServer(
			endpoints.FollowEndpoint,
			DecodeGRPCFollowRequest,
			EncodeGRPCFollowResponse,
			serverOptions...,
		),
		unfollow: grpctransport.NewServer(
			endpoints.UnFollowEndpoint,
			DecodeGRPCUnFollowRequest,
			EncodeGRPCUnFollowResponse,
			serverOptions...,
		),
		listmessage: grpctransport.NewServer(
			endpoints.ListMessageEndpoint,
			DecodeGRPCListMessageRequest,
			EncodeGRPCListMessageResponse,
			serverOptions...,
		),
		readmessage: grpctransport.NewServer(
			endpoints.ReadMessageEndpoint,
			DecodeGRPCReadMessageRequest,
			EncodeGRPCReadMessageResponse,
			serverOptions...,
		),
		getmessagesummary: grpctransport.NewServer(
			endpoints.GetMessageSummaryEndpoint,
			DecodeGRPCGetMessageSummaryRequest,
			EncodeGRPCGetMessageSummaryResponse,
			serverOptions...,
		),
		listcomment: grpctransport.NewServer(
			endpoints.ListCommentEndpoint,
			DecodeGRPCListCommentRequest,
			EncodeGRPCListCommentResponse,
			serverOptions...,
		),
		newrecommendstatus: grpctransport.NewServer(
			endpoints.NewRecommendStatusEndpoint,
			DecodeGRPCNewRecommendStatusRequest,
			EncodeGRPCNewRecommendStatusResponse,
			serverOptions...,
		),
		createcomment: grpctransport.NewServer(
			endpoints.CreateCommentEndpoint,
			DecodeGRPCCreateCommentRequest,
			EncodeGRPCCreateCommentResponse,
			serverOptions...,
		),
		likecomment: grpctransport.NewServer(
			endpoints.LikeCommentEndpoint,
			DecodeGRPCLikeCommentRequest,
			EncodeGRPCLikeCommentResponse,
			serverOptions...,
		),
		unlikecomment: grpctransport.NewServer(
			endpoints.UnlikeCommentEndpoint,
			DecodeGRPCUnlikeCommentRequest,
			EncodeGRPCUnlikeCommentResponse,
			serverOptions...,
		),
		listblacklist: grpctransport.NewServer(
			endpoints.ListBlacklistEndpoint,
			DecodeGRPCListBlacklistRequest,
			EncodeGRPCListBlacklistResponse,
			serverOptions...,
		),
		createblacklist: grpctransport.NewServer(
			endpoints.CreateBlacklistEndpoint,
			DecodeGRPCCreateBlacklistRequest,
			EncodeGRPCCreateBlacklistResponse,
			serverOptions...,
		),
		deleteblacklist: grpctransport.NewServer(
			endpoints.DeleteBlacklistEndpoint,
			DecodeGRPCDeleteBlacklistRequest,
			EncodeGRPCDeleteBlacklistResponse,
			serverOptions...,
		),
	}
}

// grpcServer implements the SocialServer interface
type grpcServer struct {
	signin             grpctransport.Handler
	finduser           grpctransport.Handler
	updateuserprofile  grpctransport.Handler
	updateuseravatar   grpctransport.Handler
	updateusername     grpctransport.Handler
	createstatus       grpctransport.Handler
	updatestatus       grpctransport.Handler
	deletestatus       grpctransport.Handler
	likestatus         grpctransport.Handler
	unlikestatus       grpctransport.Handler
	listlikestatus     grpctransport.Handler
	getstatus          grpctransport.Handler
	liststatus         grpctransport.Handler
	listrecommended    grpctransport.Handler
	listusertimeline   grpctransport.Handler
	latestfollowing    grpctransport.Handler
	listrelationship   grpctransport.Handler
	follow             grpctransport.Handler
	unfollow           grpctransport.Handler
	listmessage        grpctransport.Handler
	readmessage        grpctransport.Handler
	getmessagesummary  grpctransport.Handler
	listcomment        grpctransport.Handler
	newrecommendstatus grpctransport.Handler
	createcomment      grpctransport.Handler
	likecomment        grpctransport.Handler
	unlikecomment      grpctransport.Handler
	listblacklist      grpctransport.Handler
	createblacklist    grpctransport.Handler
	deleteblacklist    grpctransport.Handler
}

// Methods for grpcServer to implement SocialServer interface

func (s *grpcServer) SignIn(ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {
	_, rep, err := s.signin.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SignInResponse), nil
}

func (s *grpcServer) FindUser(ctx context.Context, req *pb.FindUserRequest) (*pb.FindUserResponse, error) {
	_, rep, err := s.finduser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.FindUserResponse), nil
}

func (s *grpcServer) UpdateUserProfile(ctx context.Context, req *pb.UpdateUserProfileRequest) (*pb.UpdateUserResponse, error) {
	_, rep, err := s.updateuserprofile.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateUserResponse), nil
}

func (s *grpcServer) UpdateUserAvatar(ctx context.Context, req *pb.UpdateUserAvatarRequest) (*pb.UpdateUserResponse, error) {
	_, rep, err := s.updateuseravatar.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateUserResponse), nil
}

func (s *grpcServer) UpdateUserName(ctx context.Context, req *pb.UpdateUserNameRequest) (*pb.UpdateUserResponse, error) {
	_, rep, err := s.updateusername.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateUserResponse), nil
}

func (s *grpcServer) CreateStatus(ctx context.Context, req *pb.CreateStatusRequest) (*pb.CreateStatusResponse, error) {
	_, rep, err := s.createstatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateStatusResponse), nil
}

func (s *grpcServer) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {
	_, rep, err := s.updatestatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UpdateStatusResponse), nil
}

func (s *grpcServer) DeleteStatus(ctx context.Context, req *pb.DeleteStatusRequest) (*pb.SimpleResponse, error) {
	_, rep, err := s.deletestatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SimpleResponse), nil
}

func (s *grpcServer) LikeStatus(ctx context.Context, req *pb.LikeStatusRequest) (*pb.SimpleResponse, error) {
	_, rep, err := s.likestatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SimpleResponse), nil
}

func (s *grpcServer) UnLikeStatus(ctx context.Context, req *pb.UnLikeStatusRequest) (*pb.SimpleResponse, error) {
	_, rep, err := s.unlikestatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SimpleResponse), nil
}

func (s *grpcServer) ListLikeStatus(ctx context.Context, req *pb.ListLikeRequest) (*pb.ListLikeResponse, error) {
	_, rep, err := s.listlikestatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListLikeResponse), nil
}

func (s *grpcServer) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	_, rep, err := s.getstatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetStatusResponse), nil
}

func (s *grpcServer) ListStatus(ctx context.Context, req *pb.ListStatusRequest) (*pb.ListStatusResponse, error) {
	_, rep, err := s.liststatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListStatusResponse), nil
}

func (s *grpcServer) ListRecommended(ctx context.Context, req *pb.ListStatusRequest) (*pb.ListStatusResponse, error) {
	_, rep, err := s.listrecommended.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListStatusResponse), nil
}

func (s *grpcServer) ListUserTimeline(ctx context.Context, req *pb.ListStatusRequest) (*pb.ListStatusResponse, error) {
	_, rep, err := s.listusertimeline.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListStatusResponse), nil
}

func (s *grpcServer) LatestFollowing(ctx context.Context, req *pb.LatestFollowingRequest) (*pb.LatestFollowingResponse, error) {
	_, rep, err := s.latestfollowing.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.LatestFollowingResponse), nil
}

func (s *grpcServer) ListRelationship(ctx context.Context, req *pb.ListRelationshipRequest) (*pb.ListRelationshipResponse, error) {
	_, rep, err := s.listrelationship.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListRelationshipResponse), nil
}

func (s *grpcServer) Follow(ctx context.Context, req *pb.FollowRequest) (*pb.SimpleResponse, error) {
	_, rep, err := s.follow.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SimpleResponse), nil
}

func (s *grpcServer) UnFollow(ctx context.Context, req *pb.UnFollowRequest) (*pb.SimpleResponse, error) {
	_, rep, err := s.unfollow.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SimpleResponse), nil
}

func (s *grpcServer) ListMessage(ctx context.Context, req *pb.ListMessageRequest) (*pb.ListMessageResponse, error) {
	_, rep, err := s.listmessage.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListMessageResponse), nil
}

func (s *grpcServer) ReadMessage(ctx context.Context, req *pb.ReadMessageRequest) (*pb.SimpleResponse, error) {
	_, rep, err := s.readmessage.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SimpleResponse), nil
}

func (s *grpcServer) GetMessageSummary(ctx context.Context, req *pb.GetMessageSummaryRequest) (*pb.MessageSummaryResponse, error) {
	_, rep, err := s.getmessagesummary.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MessageSummaryResponse), nil
}

func (s *grpcServer) ListComment(ctx context.Context, req *pb.ListCommentRequest) (*pb.ListCommentResponse, error) {
	_, rep, err := s.listcomment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListCommentResponse), nil
}

func (s *grpcServer) NewRecommendStatus(ctx context.Context, req *pb.NewRecommendStatusResquest) (*pb.NewRecommendStatusResponse, error) {
	_, rep, err := s.newrecommendstatus.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.NewRecommendStatusResponse), nil
}

func (s *grpcServer) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	_, rep, err := s.createcomment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateCommentResponse), nil
}

func (s *grpcServer) LikeComment(ctx context.Context, req *pb.LikeCommentRequest) (*pb.SimpleResponse, error) {
	_, rep, err := s.likecomment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SimpleResponse), nil
}

func (s *grpcServer) UnlikeComment(ctx context.Context, req *pb.UnlikeCommentRequest) (*pb.SimpleResponse, error) {
	_, rep, err := s.unlikecomment.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SimpleResponse), nil
}

func (s *grpcServer) ListBlacklist(ctx context.Context, req *pb.ListBlacklistRequest) (*pb.ListBlacklistResponse, error) {
	_, rep, err := s.listblacklist.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ListBlacklistResponse), nil
}

func (s *grpcServer) CreateBlacklist(ctx context.Context, req *pb.CreateBlacklistRequest) (*pb.SimpleResponse, error) {
	_, rep, err := s.createblacklist.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SimpleResponse), nil
}

func (s *grpcServer) DeleteBlacklist(ctx context.Context, req *pb.DeleteBlacklistRequest) (*pb.SimpleResponse, error) {
	_, rep, err := s.deleteblacklist.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.SimpleResponse), nil
}

// Server Decode

// DecodeGRPCSignInRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC signin request to a user-domain signin request. Primarily useful in a server.
func DecodeGRPCSignInRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.SignInRequest)
	return req, nil
}

// DecodeGRPCFindUserRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC finduser request to a user-domain finduser request. Primarily useful in a server.
func DecodeGRPCFindUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.FindUserRequest)
	return req, nil
}

// DecodeGRPCUpdateUserProfileRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC updateuserprofile request to a user-domain updateuserprofile request. Primarily useful in a server.
func DecodeGRPCUpdateUserProfileRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UpdateUserProfileRequest)
	return req, nil
}

// DecodeGRPCUpdateUserAvatarRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC updateuseravatar request to a user-domain updateuseravatar request. Primarily useful in a server.
func DecodeGRPCUpdateUserAvatarRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UpdateUserAvatarRequest)
	return req, nil
}

// DecodeGRPCUpdateUserNameRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC updateusername request to a user-domain updateusername request. Primarily useful in a server.
func DecodeGRPCUpdateUserNameRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UpdateUserNameRequest)
	return req, nil
}

// DecodeGRPCCreateStatusRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC createstatus request to a user-domain createstatus request. Primarily useful in a server.
func DecodeGRPCCreateStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateStatusRequest)
	return req, nil
}

// DecodeGRPCUpdateStatusRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC updatestatus request to a user-domain updatestatus request. Primarily useful in a server.
func DecodeGRPCUpdateStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UpdateStatusRequest)
	return req, nil
}

// DecodeGRPCDeleteStatusRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC deletestatus request to a user-domain deletestatus request. Primarily useful in a server.
func DecodeGRPCDeleteStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteStatusRequest)
	return req, nil
}

// DecodeGRPCLikeStatusRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC likestatus request to a user-domain likestatus request. Primarily useful in a server.
func DecodeGRPCLikeStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LikeStatusRequest)
	return req, nil
}

// DecodeGRPCUnLikeStatusRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC unlikestatus request to a user-domain unlikestatus request. Primarily useful in a server.
func DecodeGRPCUnLikeStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UnLikeStatusRequest)
	return req, nil
}

// DecodeGRPCListLikeStatusRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC listlikestatus request to a user-domain listlikestatus request. Primarily useful in a server.
func DecodeGRPCListLikeStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ListLikeRequest)
	return req, nil
}

// DecodeGRPCGetStatusRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC getstatus request to a user-domain getstatus request. Primarily useful in a server.
func DecodeGRPCGetStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetStatusRequest)
	return req, nil
}

// DecodeGRPCListStatusRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC liststatus request to a user-domain liststatus request. Primarily useful in a server.
func DecodeGRPCListStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ListStatusRequest)
	return req, nil
}

// DecodeGRPCListRecommendedRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC listrecommended request to a user-domain listrecommended request. Primarily useful in a server.
func DecodeGRPCListRecommendedRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ListStatusRequest)
	return req, nil
}

// DecodeGRPCListUserTimelineRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC listusertimeline request to a user-domain listusertimeline request. Primarily useful in a server.
func DecodeGRPCListUserTimelineRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ListStatusRequest)
	return req, nil
}

// DecodeGRPCLatestFollowingRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC latestfollowing request to a user-domain latestfollowing request. Primarily useful in a server.
func DecodeGRPCLatestFollowingRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LatestFollowingRequest)
	return req, nil
}

// DecodeGRPCListRelationshipRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC listrelationship request to a user-domain listrelationship request. Primarily useful in a server.
func DecodeGRPCListRelationshipRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ListRelationshipRequest)
	return req, nil
}

// DecodeGRPCFollowRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC follow request to a user-domain follow request. Primarily useful in a server.
func DecodeGRPCFollowRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.FollowRequest)
	return req, nil
}

// DecodeGRPCUnFollowRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC unfollow request to a user-domain unfollow request. Primarily useful in a server.
func DecodeGRPCUnFollowRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UnFollowRequest)
	return req, nil
}

// DecodeGRPCListMessageRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC listmessage request to a user-domain listmessage request. Primarily useful in a server.
func DecodeGRPCListMessageRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ListMessageRequest)
	return req, nil
}

// DecodeGRPCReadMessageRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC readmessage request to a user-domain readmessage request. Primarily useful in a server.
func DecodeGRPCReadMessageRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ReadMessageRequest)
	return req, nil
}

// DecodeGRPCGetMessageSummaryRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC getmessagesummary request to a user-domain getmessagesummary request. Primarily useful in a server.
func DecodeGRPCGetMessageSummaryRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetMessageSummaryRequest)
	return req, nil
}

// DecodeGRPCListCommentRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC listcomment request to a user-domain listcomment request. Primarily useful in a server.
func DecodeGRPCListCommentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ListCommentRequest)
	return req, nil
}

// DecodeGRPCNewRecommendStatusRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC newrecommendstatus request to a user-domain newrecommendstatus request. Primarily useful in a server.
func DecodeGRPCNewRecommendStatusRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.NewRecommendStatusResquest)
	return req, nil
}

// DecodeGRPCCreateCommentRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC createcomment request to a user-domain createcomment request. Primarily useful in a server.
func DecodeGRPCCreateCommentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateCommentRequest)
	return req, nil
}

// DecodeGRPCLikeCommentRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC likecomment request to a user-domain likecomment request. Primarily useful in a server.
func DecodeGRPCLikeCommentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LikeCommentRequest)
	return req, nil
}

// DecodeGRPCUnlikeCommentRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC unlikecomment request to a user-domain unlikecomment request. Primarily useful in a server.
func DecodeGRPCUnlikeCommentRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UnlikeCommentRequest)
	return req, nil
}

// DecodeGRPCListBlacklistRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC listblacklist request to a user-domain listblacklist request. Primarily useful in a server.
func DecodeGRPCListBlacklistRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ListBlacklistRequest)
	return req, nil
}

// DecodeGRPCCreateBlacklistRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC createblacklist request to a user-domain createblacklist request. Primarily useful in a server.
func DecodeGRPCCreateBlacklistRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateBlacklistRequest)
	return req, nil
}

// DecodeGRPCDeleteBlacklistRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC deleteblacklist request to a user-domain deleteblacklist request. Primarily useful in a server.
func DecodeGRPCDeleteBlacklistRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeleteBlacklistRequest)
	return req, nil
}

// Server Encode

// EncodeGRPCSignInResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain signin response to a gRPC signin reply. Primarily useful in a server.
func EncodeGRPCSignInResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SignInResponse)
	return resp, nil
}

// EncodeGRPCFindUserResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain finduser response to a gRPC finduser reply. Primarily useful in a server.
func EncodeGRPCFindUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.FindUserResponse)
	return resp, nil
}

// EncodeGRPCUpdateUserProfileResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain updateuserprofile response to a gRPC updateuserprofile reply. Primarily useful in a server.
func EncodeGRPCUpdateUserProfileResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UpdateUserResponse)
	return resp, nil
}

// EncodeGRPCUpdateUserAvatarResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain updateuseravatar response to a gRPC updateuseravatar reply. Primarily useful in a server.
func EncodeGRPCUpdateUserAvatarResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UpdateUserResponse)
	return resp, nil
}

// EncodeGRPCUpdateUserNameResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain updateusername response to a gRPC updateusername reply. Primarily useful in a server.
func EncodeGRPCUpdateUserNameResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UpdateUserResponse)
	return resp, nil
}

// EncodeGRPCCreateStatusResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain createstatus response to a gRPC createstatus reply. Primarily useful in a server.
func EncodeGRPCCreateStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.CreateStatusResponse)
	return resp, nil
}

// EncodeGRPCUpdateStatusResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain updatestatus response to a gRPC updatestatus reply. Primarily useful in a server.
func EncodeGRPCUpdateStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UpdateStatusResponse)
	return resp, nil
}

// EncodeGRPCDeleteStatusResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain deletestatus response to a gRPC deletestatus reply. Primarily useful in a server.
func EncodeGRPCDeleteStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SimpleResponse)
	return resp, nil
}

// EncodeGRPCLikeStatusResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain likestatus response to a gRPC likestatus reply. Primarily useful in a server.
func EncodeGRPCLikeStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SimpleResponse)
	return resp, nil
}

// EncodeGRPCUnLikeStatusResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain unlikestatus response to a gRPC unlikestatus reply. Primarily useful in a server.
func EncodeGRPCUnLikeStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SimpleResponse)
	return resp, nil
}

// EncodeGRPCListLikeStatusResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain listlikestatus response to a gRPC listlikestatus reply. Primarily useful in a server.
func EncodeGRPCListLikeStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ListLikeResponse)
	return resp, nil
}

// EncodeGRPCGetStatusResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain getstatus response to a gRPC getstatus reply. Primarily useful in a server.
func EncodeGRPCGetStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.GetStatusResponse)
	return resp, nil
}

// EncodeGRPCListStatusResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain liststatus response to a gRPC liststatus reply. Primarily useful in a server.
func EncodeGRPCListStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ListStatusResponse)
	return resp, nil
}

// EncodeGRPCListRecommendedResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain listrecommended response to a gRPC listrecommended reply. Primarily useful in a server.
func EncodeGRPCListRecommendedResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ListStatusResponse)
	return resp, nil
}

// EncodeGRPCListUserTimelineResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain listusertimeline response to a gRPC listusertimeline reply. Primarily useful in a server.
func EncodeGRPCListUserTimelineResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ListStatusResponse)
	return resp, nil
}

// EncodeGRPCLatestFollowingResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain latestfollowing response to a gRPC latestfollowing reply. Primarily useful in a server.
func EncodeGRPCLatestFollowingResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.LatestFollowingResponse)
	return resp, nil
}

// EncodeGRPCListRelationshipResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain listrelationship response to a gRPC listrelationship reply. Primarily useful in a server.
func EncodeGRPCListRelationshipResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ListRelationshipResponse)
	return resp, nil
}

// EncodeGRPCFollowResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain follow response to a gRPC follow reply. Primarily useful in a server.
func EncodeGRPCFollowResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SimpleResponse)
	return resp, nil
}

// EncodeGRPCUnFollowResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain unfollow response to a gRPC unfollow reply. Primarily useful in a server.
func EncodeGRPCUnFollowResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SimpleResponse)
	return resp, nil
}

// EncodeGRPCListMessageResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain listmessage response to a gRPC listmessage reply. Primarily useful in a server.
func EncodeGRPCListMessageResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ListMessageResponse)
	return resp, nil
}

// EncodeGRPCReadMessageResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain readmessage response to a gRPC readmessage reply. Primarily useful in a server.
func EncodeGRPCReadMessageResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SimpleResponse)
	return resp, nil
}

// EncodeGRPCGetMessageSummaryResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain getmessagesummary response to a gRPC getmessagesummary reply. Primarily useful in a server.
func EncodeGRPCGetMessageSummaryResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.MessageSummaryResponse)
	return resp, nil
}

// EncodeGRPCListCommentResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain listcomment response to a gRPC listcomment reply. Primarily useful in a server.
func EncodeGRPCListCommentResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ListCommentResponse)
	return resp, nil
}

// EncodeGRPCNewRecommendStatusResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain newrecommendstatus response to a gRPC newrecommendstatus reply. Primarily useful in a server.
func EncodeGRPCNewRecommendStatusResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.NewRecommendStatusResponse)
	return resp, nil
}

// EncodeGRPCCreateCommentResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain createcomment response to a gRPC createcomment reply. Primarily useful in a server.
func EncodeGRPCCreateCommentResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.CreateCommentResponse)
	return resp, nil
}

// EncodeGRPCLikeCommentResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain likecomment response to a gRPC likecomment reply. Primarily useful in a server.
func EncodeGRPCLikeCommentResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SimpleResponse)
	return resp, nil
}

// EncodeGRPCUnlikeCommentResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain unlikecomment response to a gRPC unlikecomment reply. Primarily useful in a server.
func EncodeGRPCUnlikeCommentResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SimpleResponse)
	return resp, nil
}

// EncodeGRPCListBlacklistResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain listblacklist response to a gRPC listblacklist reply. Primarily useful in a server.
func EncodeGRPCListBlacklistResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ListBlacklistResponse)
	return resp, nil
}

// EncodeGRPCCreateBlacklistResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain createblacklist response to a gRPC createblacklist reply. Primarily useful in a server.
func EncodeGRPCCreateBlacklistResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SimpleResponse)
	return resp, nil
}

// EncodeGRPCDeleteBlacklistResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain deleteblacklist response to a gRPC deleteblacklist reply. Primarily useful in a server.
func EncodeGRPCDeleteBlacklistResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.SimpleResponse)
	return resp, nil
}

// Helpers

func metadataToContext(ctx context.Context, md metadata.MD) context.Context {
	for k, v := range md {
		if v != nil {
			// The key is added both in metadata format (k) which is all lower
			// and the http.CanonicalHeaderKey of the key so that it can be
			// accessed in either format
			ctx = context.WithValue(ctx, k, v[0])
			ctx = context.WithValue(ctx, http.CanonicalHeaderKey(k), v[0])
		}
	}

	return ctx
}
