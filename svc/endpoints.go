// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: fcd9ff140d
// Version Date: 2021-07-14T06:36:40Z

package svc

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints, as well as encoders and
// decoders for those types, for all of our supported transport serialization
// formats.

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	pb "github.com/mises-id/socialsvc/proto"
)

// Endpoints collects all of the endpoints that compose an add service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	SignInEndpoint            endpoint.Endpoint
	FindUserEndpoint          endpoint.Endpoint
	UpdateUserProfileEndpoint endpoint.Endpoint
	UpdateUserAvatarEndpoint  endpoint.Endpoint
	UpdateUserNameEndpoint    endpoint.Endpoint
	CreateStatusEndpoint      endpoint.Endpoint
	DeleteStatusEndpoint      endpoint.Endpoint
	LikeStatusEndpoint        endpoint.Endpoint
	UnLikeStatusEndpoint      endpoint.Endpoint
	GetStatusEndpoint         endpoint.Endpoint
	ListStatusEndpoint        endpoint.Endpoint
	ListRecommendedEndpoint   endpoint.Endpoint
	ListUserTimelineEndpoint  endpoint.Endpoint
	ListRelationshipEndpoint  endpoint.Endpoint
	FollowEndpoint            endpoint.Endpoint
	UnFollowEndpoint          endpoint.Endpoint
}

// Endpoints

func (e Endpoints) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	response, err := e.SignInEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.SignInResponse), nil
}

func (e Endpoints) FindUser(ctx context.Context, in *pb.FindUserRequest) (*pb.FindUserResponse, error) {
	response, err := e.FindUserEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.FindUserResponse), nil
}

func (e Endpoints) UpdateUserProfile(ctx context.Context, in *pb.UpdateUserProfileRequest) (*pb.UpdateUserResponse, error) {
	response, err := e.UpdateUserProfileEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UpdateUserResponse), nil
}

func (e Endpoints) UpdateUserAvatar(ctx context.Context, in *pb.UpdateUserAvatarRequest) (*pb.UpdateUserResponse, error) {
	response, err := e.UpdateUserAvatarEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UpdateUserResponse), nil
}

func (e Endpoints) UpdateUserName(ctx context.Context, in *pb.UpdateUserNameRequest) (*pb.UpdateUserResponse, error) {
	response, err := e.UpdateUserNameEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.UpdateUserResponse), nil
}

func (e Endpoints) CreateStatus(ctx context.Context, in *pb.CreateStatusRequest) (*pb.CreateStatusResponse, error) {
	response, err := e.CreateStatusEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.CreateStatusResponse), nil
}

func (e Endpoints) DeleteStatus(ctx context.Context, in *pb.DeleteStatusRequest) (*pb.SimpleResponse, error) {
	response, err := e.DeleteStatusEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.SimpleResponse), nil
}

func (e Endpoints) LikeStatus(ctx context.Context, in *pb.LikeStatusRequest) (*pb.SimpleResponse, error) {
	response, err := e.LikeStatusEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.SimpleResponse), nil
}

func (e Endpoints) UnLikeStatus(ctx context.Context, in *pb.UnLikeStatusRequest) (*pb.SimpleResponse, error) {
	response, err := e.UnLikeStatusEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.SimpleResponse), nil
}

func (e Endpoints) GetStatus(ctx context.Context, in *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	response, err := e.GetStatusEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.GetStatusResponse), nil
}

func (e Endpoints) ListStatus(ctx context.Context, in *pb.ListStatusRequest) (*pb.ListStatusResponse, error) {
	response, err := e.ListStatusEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.ListStatusResponse), nil
}

func (e Endpoints) ListRecommended(ctx context.Context, in *pb.ListStatusRequest) (*pb.ListStatusResponse, error) {
	response, err := e.ListRecommendedEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.ListStatusResponse), nil
}

func (e Endpoints) ListUserTimeline(ctx context.Context, in *pb.ListStatusRequest) (*pb.ListStatusResponse, error) {
	response, err := e.ListUserTimelineEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.ListStatusResponse), nil
}

func (e Endpoints) ListRelationship(ctx context.Context, in *pb.ListRelationshipRequest) (*pb.ListRelationshipResponse, error) {
	response, err := e.ListRelationshipEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.ListRelationshipResponse), nil
}

func (e Endpoints) Follow(ctx context.Context, in *pb.FollowRequest) (*pb.SimpleResponse, error) {
	response, err := e.FollowEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.SimpleResponse), nil
}

func (e Endpoints) UnFollow(ctx context.Context, in *pb.UnFollowRequest) (*pb.SimpleResponse, error) {
	response, err := e.UnFollowEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.SimpleResponse), nil
}

// Make Endpoints

func MakeSignInEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.SignInRequest)
		v, err := s.SignIn(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeFindUserEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.FindUserRequest)
		v, err := s.FindUser(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeUpdateUserProfileEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UpdateUserProfileRequest)
		v, err := s.UpdateUserProfile(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeUpdateUserAvatarEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UpdateUserAvatarRequest)
		v, err := s.UpdateUserAvatar(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeUpdateUserNameEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UpdateUserNameRequest)
		v, err := s.UpdateUserName(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeCreateStatusEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.CreateStatusRequest)
		v, err := s.CreateStatus(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeDeleteStatusEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.DeleteStatusRequest)
		v, err := s.DeleteStatus(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeLikeStatusEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.LikeStatusRequest)
		v, err := s.LikeStatus(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeUnLikeStatusEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UnLikeStatusRequest)
		v, err := s.UnLikeStatus(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeGetStatusEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.GetStatusRequest)
		v, err := s.GetStatus(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeListStatusEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.ListStatusRequest)
		v, err := s.ListStatus(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeListRecommendedEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.ListStatusRequest)
		v, err := s.ListRecommended(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeListUserTimelineEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.ListStatusRequest)
		v, err := s.ListUserTimeline(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeListRelationshipEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.ListRelationshipRequest)
		v, err := s.ListRelationship(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeFollowEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.FollowRequest)
		v, err := s.Follow(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeUnFollowEndpoint(s pb.SocialServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.UnFollowRequest)
		v, err := s.UnFollow(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// WrapAllExcept wraps each Endpoint field of struct Endpoints with a
// go-kit/kit/endpoint.Middleware.
// Use this for applying a set of middlewares to every endpoint in the service.
// Optionally, endpoints can be passed in by name to be excluded from being wrapped.
// WrapAllExcept(middleware, "Status", "Ping")
func (e *Endpoints) WrapAllExcept(middleware endpoint.Middleware, excluded ...string) {
	included := map[string]struct{}{
		"SignIn":            {},
		"FindUser":          {},
		"UpdateUserProfile": {},
		"UpdateUserAvatar":  {},
		"UpdateUserName":    {},
		"CreateStatus":      {},
		"DeleteStatus":      {},
		"LikeStatus":        {},
		"UnLikeStatus":      {},
		"GetStatus":         {},
		"ListStatus":        {},
		"ListRecommended":   {},
		"ListUserTimeline":  {},
		"ListRelationship":  {},
		"Follow":            {},
		"UnFollow":          {},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc := range included {
		if inc == "SignIn" {
			e.SignInEndpoint = middleware(e.SignInEndpoint)
		}
		if inc == "FindUser" {
			e.FindUserEndpoint = middleware(e.FindUserEndpoint)
		}
		if inc == "UpdateUserProfile" {
			e.UpdateUserProfileEndpoint = middleware(e.UpdateUserProfileEndpoint)
		}
		if inc == "UpdateUserAvatar" {
			e.UpdateUserAvatarEndpoint = middleware(e.UpdateUserAvatarEndpoint)
		}
		if inc == "UpdateUserName" {
			e.UpdateUserNameEndpoint = middleware(e.UpdateUserNameEndpoint)
		}
		if inc == "CreateStatus" {
			e.CreateStatusEndpoint = middleware(e.CreateStatusEndpoint)
		}
		if inc == "DeleteStatus" {
			e.DeleteStatusEndpoint = middleware(e.DeleteStatusEndpoint)
		}
		if inc == "LikeStatus" {
			e.LikeStatusEndpoint = middleware(e.LikeStatusEndpoint)
		}
		if inc == "UnLikeStatus" {
			e.UnLikeStatusEndpoint = middleware(e.UnLikeStatusEndpoint)
		}
		if inc == "GetStatus" {
			e.GetStatusEndpoint = middleware(e.GetStatusEndpoint)
		}
		if inc == "ListStatus" {
			e.ListStatusEndpoint = middleware(e.ListStatusEndpoint)
		}
		if inc == "ListRecommended" {
			e.ListRecommendedEndpoint = middleware(e.ListRecommendedEndpoint)
		}
		if inc == "ListUserTimeline" {
			e.ListUserTimelineEndpoint = middleware(e.ListUserTimelineEndpoint)
		}
		if inc == "ListRelationship" {
			e.ListRelationshipEndpoint = middleware(e.ListRelationshipEndpoint)
		}
		if inc == "Follow" {
			e.FollowEndpoint = middleware(e.FollowEndpoint)
		}
		if inc == "UnFollow" {
			e.UnFollowEndpoint = middleware(e.UnFollowEndpoint)
		}
	}
}

// LabeledMiddleware will get passed the endpoint name when passed to
// WrapAllLabeledExcept, this can be used to write a generic metrics
// middleware which can send the endpoint name to the metrics collector.
type LabeledMiddleware func(string, endpoint.Endpoint) endpoint.Endpoint

// WrapAllLabeledExcept wraps each Endpoint field of struct Endpoints with a
// LabeledMiddleware, which will receive the name of the endpoint. See
// LabeldMiddleware. See method WrapAllExept for details on excluded
// functionality.
func (e *Endpoints) WrapAllLabeledExcept(middleware func(string, endpoint.Endpoint) endpoint.Endpoint, excluded ...string) {
	included := map[string]struct{}{
		"SignIn":            {},
		"FindUser":          {},
		"UpdateUserProfile": {},
		"UpdateUserAvatar":  {},
		"UpdateUserName":    {},
		"CreateStatus":      {},
		"DeleteStatus":      {},
		"LikeStatus":        {},
		"UnLikeStatus":      {},
		"GetStatus":         {},
		"ListStatus":        {},
		"ListRecommended":   {},
		"ListUserTimeline":  {},
		"ListRelationship":  {},
		"Follow":            {},
		"UnFollow":          {},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc := range included {
		if inc == "SignIn" {
			e.SignInEndpoint = middleware("SignIn", e.SignInEndpoint)
		}
		if inc == "FindUser" {
			e.FindUserEndpoint = middleware("FindUser", e.FindUserEndpoint)
		}
		if inc == "UpdateUserProfile" {
			e.UpdateUserProfileEndpoint = middleware("UpdateUserProfile", e.UpdateUserProfileEndpoint)
		}
		if inc == "UpdateUserAvatar" {
			e.UpdateUserAvatarEndpoint = middleware("UpdateUserAvatar", e.UpdateUserAvatarEndpoint)
		}
		if inc == "UpdateUserName" {
			e.UpdateUserNameEndpoint = middleware("UpdateUserName", e.UpdateUserNameEndpoint)
		}
		if inc == "CreateStatus" {
			e.CreateStatusEndpoint = middleware("CreateStatus", e.CreateStatusEndpoint)
		}
		if inc == "DeleteStatus" {
			e.DeleteStatusEndpoint = middleware("DeleteStatus", e.DeleteStatusEndpoint)
		}
		if inc == "LikeStatus" {
			e.LikeStatusEndpoint = middleware("LikeStatus", e.LikeStatusEndpoint)
		}
		if inc == "UnLikeStatus" {
			e.UnLikeStatusEndpoint = middleware("UnLikeStatus", e.UnLikeStatusEndpoint)
		}
		if inc == "GetStatus" {
			e.GetStatusEndpoint = middleware("GetStatus", e.GetStatusEndpoint)
		}
		if inc == "ListStatus" {
			e.ListStatusEndpoint = middleware("ListStatus", e.ListStatusEndpoint)
		}
		if inc == "ListRecommended" {
			e.ListRecommendedEndpoint = middleware("ListRecommended", e.ListRecommendedEndpoint)
		}
		if inc == "ListUserTimeline" {
			e.ListUserTimelineEndpoint = middleware("ListUserTimeline", e.ListUserTimelineEndpoint)
		}
		if inc == "ListRelationship" {
			e.ListRelationshipEndpoint = middleware("ListRelationship", e.ListRelationshipEndpoint)
		}
		if inc == "Follow" {
			e.FollowEndpoint = middleware("Follow", e.FollowEndpoint)
		}
		if inc == "UnFollow" {
			e.UnFollowEndpoint = middleware("UnFollow", e.UnFollowEndpoint)
		}
	}
}
