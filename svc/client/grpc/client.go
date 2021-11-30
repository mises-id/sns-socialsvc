// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: fcd9ff140d
// Version Date: 2021-07-14T06:36:40Z

// Package grpc provides a gRPC client for the Social service.
package grpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/mises-id/socialsvc/proto"
	"github.com/mises-id/socialsvc/svc"
)

// New returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn, options ...ClientOption) (pb.SocialServer, error) {
	var cc clientConfig

	for _, f := range options {
		err := f(&cc)
		if err != nil {
			return nil, errors.Wrap(err, "cannot apply option")
		}
	}

	clientOptions := []grpctransport.ClientOption{
		grpctransport.ClientBefore(
			contextValuesToGRPCMetadata(cc.headers)),
	}
	var signinEndpoint endpoint.Endpoint
	{
		signinEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"SignIn",
			EncodeGRPCSignInRequest,
			DecodeGRPCSignInResponse,
			pb.SignInResponse{},
			clientOptions...,
		).Endpoint()
	}

	var finduserEndpoint endpoint.Endpoint
	{
		finduserEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"FindUser",
			EncodeGRPCFindUserRequest,
			DecodeGRPCFindUserResponse,
			pb.FindUserResponse{},
			clientOptions...,
		).Endpoint()
	}

	var updateuserprofileEndpoint endpoint.Endpoint
	{
		updateuserprofileEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"UpdateUserProfile",
			EncodeGRPCUpdateUserProfileRequest,
			DecodeGRPCUpdateUserProfileResponse,
			pb.UpdateUserResponse{},
			clientOptions...,
		).Endpoint()
	}

	var updateuseravatarEndpoint endpoint.Endpoint
	{
		updateuseravatarEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"UpdateUserAvatar",
			EncodeGRPCUpdateUserAvatarRequest,
			DecodeGRPCUpdateUserAvatarResponse,
			pb.UpdateUserResponse{},
			clientOptions...,
		).Endpoint()
	}

	var updateusernameEndpoint endpoint.Endpoint
	{
		updateusernameEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"UpdateUserName",
			EncodeGRPCUpdateUserNameRequest,
			DecodeGRPCUpdateUserNameResponse,
			pb.UpdateUserResponse{},
			clientOptions...,
		).Endpoint()
	}

	var createstatusEndpoint endpoint.Endpoint
	{
		createstatusEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"CreateStatus",
			EncodeGRPCCreateStatusRequest,
			DecodeGRPCCreateStatusResponse,
			pb.CreateStatusResponse{},
			clientOptions...,
		).Endpoint()
	}

	var deletestatusEndpoint endpoint.Endpoint
	{
		deletestatusEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"DeleteStatus",
			EncodeGRPCDeleteStatusRequest,
			DecodeGRPCDeleteStatusResponse,
			pb.SimpleResponse{},
			clientOptions...,
		).Endpoint()
	}

	var likestatusEndpoint endpoint.Endpoint
	{
		likestatusEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"LikeStatus",
			EncodeGRPCLikeStatusRequest,
			DecodeGRPCLikeStatusResponse,
			pb.SimpleResponse{},
			clientOptions...,
		).Endpoint()
	}

	var unlikestatusEndpoint endpoint.Endpoint
	{
		unlikestatusEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"UnLikeStatus",
			EncodeGRPCUnLikeStatusRequest,
			DecodeGRPCUnLikeStatusResponse,
			pb.SimpleResponse{},
			clientOptions...,
		).Endpoint()
	}

	var getstatusEndpoint endpoint.Endpoint
	{
		getstatusEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"GetStatus",
			EncodeGRPCGetStatusRequest,
			DecodeGRPCGetStatusResponse,
			pb.GetStatusResponse{},
			clientOptions...,
		).Endpoint()
	}

	var liststatusEndpoint endpoint.Endpoint
	{
		liststatusEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"ListStatus",
			EncodeGRPCListStatusRequest,
			DecodeGRPCListStatusResponse,
			pb.ListStatusResponse{},
			clientOptions...,
		).Endpoint()
	}

	var listrecommendedEndpoint endpoint.Endpoint
	{
		listrecommendedEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"ListRecommended",
			EncodeGRPCListRecommendedRequest,
			DecodeGRPCListRecommendedResponse,
			pb.ListStatusResponse{},
			clientOptions...,
		).Endpoint()
	}

	var listusertimelineEndpoint endpoint.Endpoint
	{
		listusertimelineEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"ListUserTimeline",
			EncodeGRPCListUserTimelineRequest,
			DecodeGRPCListUserTimelineResponse,
			pb.ListStatusResponse{},
			clientOptions...,
		).Endpoint()
	}

	var listrelationshipEndpoint endpoint.Endpoint
	{
		listrelationshipEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"ListRelationship",
			EncodeGRPCListRelationshipRequest,
			DecodeGRPCListRelationshipResponse,
			pb.ListRelationshipResponse{},
			clientOptions...,
		).Endpoint()
	}

	var followEndpoint endpoint.Endpoint
	{
		followEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"Follow",
			EncodeGRPCFollowRequest,
			DecodeGRPCFollowResponse,
			pb.SimpleResponse{},
			clientOptions...,
		).Endpoint()
	}

	var unfollowEndpoint endpoint.Endpoint
	{
		unfollowEndpoint = grpctransport.NewClient(
			conn,
			"socialsvc.Social",
			"UnFollow",
			EncodeGRPCUnFollowRequest,
			DecodeGRPCUnFollowResponse,
			pb.SimpleResponse{},
			clientOptions...,
		).Endpoint()
	}

	return svc.Endpoints{
		SignInEndpoint:            signinEndpoint,
		FindUserEndpoint:          finduserEndpoint,
		UpdateUserProfileEndpoint: updateuserprofileEndpoint,
		UpdateUserAvatarEndpoint:  updateuseravatarEndpoint,
		UpdateUserNameEndpoint:    updateusernameEndpoint,
		CreateStatusEndpoint:      createstatusEndpoint,
		DeleteStatusEndpoint:      deletestatusEndpoint,
		LikeStatusEndpoint:        likestatusEndpoint,
		UnLikeStatusEndpoint:      unlikestatusEndpoint,
		GetStatusEndpoint:         getstatusEndpoint,
		ListStatusEndpoint:        liststatusEndpoint,
		ListRecommendedEndpoint:   listrecommendedEndpoint,
		ListUserTimelineEndpoint:  listusertimelineEndpoint,
		ListRelationshipEndpoint:  listrelationshipEndpoint,
		FollowEndpoint:            followEndpoint,
		UnFollowEndpoint:          unfollowEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCSignInResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC signin reply to a user-domain signin response. Primarily useful in a client.
func DecodeGRPCSignInResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SignInResponse)
	return reply, nil
}

// DecodeGRPCFindUserResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC finduser reply to a user-domain finduser response. Primarily useful in a client.
func DecodeGRPCFindUserResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.FindUserResponse)
	return reply, nil
}

// DecodeGRPCUpdateUserProfileResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC updateuserprofile reply to a user-domain updateuserprofile response. Primarily useful in a client.
func DecodeGRPCUpdateUserProfileResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UpdateUserResponse)
	return reply, nil
}

// DecodeGRPCUpdateUserAvatarResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC updateuseravatar reply to a user-domain updateuseravatar response. Primarily useful in a client.
func DecodeGRPCUpdateUserAvatarResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UpdateUserResponse)
	return reply, nil
}

// DecodeGRPCUpdateUserNameResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC updateusername reply to a user-domain updateusername response. Primarily useful in a client.
func DecodeGRPCUpdateUserNameResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UpdateUserResponse)
	return reply, nil
}

// DecodeGRPCCreateStatusResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC createstatus reply to a user-domain createstatus response. Primarily useful in a client.
func DecodeGRPCCreateStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.CreateStatusResponse)
	return reply, nil
}

// DecodeGRPCDeleteStatusResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC deletestatus reply to a user-domain deletestatus response. Primarily useful in a client.
func DecodeGRPCDeleteStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SimpleResponse)
	return reply, nil
}

// DecodeGRPCLikeStatusResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC likestatus reply to a user-domain likestatus response. Primarily useful in a client.
func DecodeGRPCLikeStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SimpleResponse)
	return reply, nil
}

// DecodeGRPCUnLikeStatusResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC unlikestatus reply to a user-domain unlikestatus response. Primarily useful in a client.
func DecodeGRPCUnLikeStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SimpleResponse)
	return reply, nil
}

// DecodeGRPCGetStatusResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC getstatus reply to a user-domain getstatus response. Primarily useful in a client.
func DecodeGRPCGetStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetStatusResponse)
	return reply, nil
}

// DecodeGRPCListStatusResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC liststatus reply to a user-domain liststatus response. Primarily useful in a client.
func DecodeGRPCListStatusResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ListStatusResponse)
	return reply, nil
}

// DecodeGRPCListRecommendedResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC listrecommended reply to a user-domain listrecommended response. Primarily useful in a client.
func DecodeGRPCListRecommendedResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ListStatusResponse)
	return reply, nil
}

// DecodeGRPCListUserTimelineResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC listusertimeline reply to a user-domain listusertimeline response. Primarily useful in a client.
func DecodeGRPCListUserTimelineResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ListStatusResponse)
	return reply, nil
}

// DecodeGRPCListRelationshipResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC listrelationship reply to a user-domain listrelationship response. Primarily useful in a client.
func DecodeGRPCListRelationshipResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ListRelationshipResponse)
	return reply, nil
}

// DecodeGRPCFollowResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC follow reply to a user-domain follow response. Primarily useful in a client.
func DecodeGRPCFollowResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SimpleResponse)
	return reply, nil
}

// DecodeGRPCUnFollowResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC unfollow reply to a user-domain unfollow response. Primarily useful in a client.
func DecodeGRPCUnFollowResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.SimpleResponse)
	return reply, nil
}

// GRPC Client Encode

// EncodeGRPCSignInRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain signin request to a gRPC signin request. Primarily useful in a client.
func EncodeGRPCSignInRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.SignInRequest)
	return req, nil
}

// EncodeGRPCFindUserRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain finduser request to a gRPC finduser request. Primarily useful in a client.
func EncodeGRPCFindUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.FindUserRequest)
	return req, nil
}

// EncodeGRPCUpdateUserProfileRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain updateuserprofile request to a gRPC updateuserprofile request. Primarily useful in a client.
func EncodeGRPCUpdateUserProfileRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UpdateUserProfileRequest)
	return req, nil
}

// EncodeGRPCUpdateUserAvatarRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain updateuseravatar request to a gRPC updateuseravatar request. Primarily useful in a client.
func EncodeGRPCUpdateUserAvatarRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UpdateUserAvatarRequest)
	return req, nil
}

// EncodeGRPCUpdateUserNameRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain updateusername request to a gRPC updateusername request. Primarily useful in a client.
func EncodeGRPCUpdateUserNameRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UpdateUserNameRequest)
	return req, nil
}

// EncodeGRPCCreateStatusRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain createstatus request to a gRPC createstatus request. Primarily useful in a client.
func EncodeGRPCCreateStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.CreateStatusRequest)
	return req, nil
}

// EncodeGRPCDeleteStatusRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain deletestatus request to a gRPC deletestatus request. Primarily useful in a client.
func EncodeGRPCDeleteStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.DeleteStatusRequest)
	return req, nil
}

// EncodeGRPCLikeStatusRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain likestatus request to a gRPC likestatus request. Primarily useful in a client.
func EncodeGRPCLikeStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.LikeStatusRequest)
	return req, nil
}

// EncodeGRPCUnLikeStatusRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain unlikestatus request to a gRPC unlikestatus request. Primarily useful in a client.
func EncodeGRPCUnLikeStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UnLikeStatusRequest)
	return req, nil
}

// EncodeGRPCGetStatusRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain getstatus request to a gRPC getstatus request. Primarily useful in a client.
func EncodeGRPCGetStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetStatusRequest)
	return req, nil
}

// EncodeGRPCListStatusRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain liststatus request to a gRPC liststatus request. Primarily useful in a client.
func EncodeGRPCListStatusRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ListStatusRequest)
	return req, nil
}

// EncodeGRPCListRecommendedRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain listrecommended request to a gRPC listrecommended request. Primarily useful in a client.
func EncodeGRPCListRecommendedRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ListStatusRequest)
	return req, nil
}

// EncodeGRPCListUserTimelineRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain listusertimeline request to a gRPC listusertimeline request. Primarily useful in a client.
func EncodeGRPCListUserTimelineRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ListStatusRequest)
	return req, nil
}

// EncodeGRPCListRelationshipRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain listrelationship request to a gRPC listrelationship request. Primarily useful in a client.
func EncodeGRPCListRelationshipRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ListRelationshipRequest)
	return req, nil
}

// EncodeGRPCFollowRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain follow request to a gRPC follow request. Primarily useful in a client.
func EncodeGRPCFollowRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.FollowRequest)
	return req, nil
}

// EncodeGRPCUnFollowRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain unfollow request to a gRPC unfollow request. Primarily useful in a client.
func EncodeGRPCUnFollowRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UnFollowRequest)
	return req, nil
}

type clientConfig struct {
	headers []string
}

// ClientOption is a function that modifies the client config
type ClientOption func(*clientConfig) error

func CtxValuesToSend(keys ...string) ClientOption {
	return func(o *clientConfig) error {
		o.headers = keys
		return nil
	}
}

func contextValuesToGRPCMetadata(keys []string) grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		var pairs []string
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				pairs = append(pairs, k, v)
			}
		}

		if pairs != nil {
			*md = metadata.Join(*md, metadata.Pairs(pairs...))
		}

		return ctx
	}
}
