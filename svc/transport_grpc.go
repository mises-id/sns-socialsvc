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
	}
}

// grpcServer implements the SocialServer interface
type grpcServer struct {
	signin            grpctransport.Handler
	finduser          grpctransport.Handler
	updateuserprofile grpctransport.Handler
	updateuseravatar  grpctransport.Handler
	updateusername    grpctransport.Handler
	createstatus      grpctransport.Handler
	deletestatus      grpctransport.Handler
	likestatus        grpctransport.Handler
	unlikestatus      grpctransport.Handler
	getstatus         grpctransport.Handler
	liststatus        grpctransport.Handler
	listrecommended   grpctransport.Handler
	listusertimeline  grpctransport.Handler
	listrelationship  grpctransport.Handler
	follow            grpctransport.Handler
	unfollow          grpctransport.Handler
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
