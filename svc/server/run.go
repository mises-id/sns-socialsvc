// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 5f7d5bf015
// Version Date: 2021-11-26T09:27:01Z

package server

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"

	// 3d Party
	"google.golang.org/grpc"

	// This Service
	"github.com/mises-id/sns-socialsvc/handlers"
	pb "github.com/mises-id/sns-socialsvc/proto"
	"github.com/mises-id/sns-socialsvc/svc"
)

var DefaultConfig svc.Config

func init() {
	flag.StringVar(&DefaultConfig.DebugAddr, "debug.addr", ":5060", "Debug and metrics listen address")
	flag.StringVar(&DefaultConfig.HTTPAddr, "http.addr", ":5050", "HTTP listen address")
	flag.StringVar(&DefaultConfig.GRPCAddr, "grpc.addr", ":5040", "gRPC (HTTP) listen address")

	// Use environment variables, if set. Flags have priority over Env vars.
	if addr := os.Getenv("DEBUG_ADDR"); addr != "" {
		DefaultConfig.DebugAddr = addr
	}
	if port := os.Getenv("PORT"); port != "" {
		DefaultConfig.HTTPAddr = fmt.Sprintf(":%s", port)
	}
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		DefaultConfig.HTTPAddr = addr
	}
	if addr := os.Getenv("GRPC_ADDR"); addr != "" {
		DefaultConfig.GRPCAddr = addr
	}
}

func NewEndpoints(service pb.SocialServer) svc.Endpoints {
	// Business domain.

	// Wrap Service with middlewares. See handlers/middlewares.go
	service = handlers.WrapService(service)

	// Endpoint domain.
	var (
		signinEndpoint            = svc.MakeSignInEndpoint(service)
		finduserEndpoint          = svc.MakeFindUserEndpoint(service)
		updateuserprofileEndpoint = svc.MakeUpdateUserProfileEndpoint(service)
		updateuseravatarEndpoint  = svc.MakeUpdateUserAvatarEndpoint(service)
		updateusernameEndpoint    = svc.MakeUpdateUserNameEndpoint(service)
		createstatusEndpoint      = svc.MakeCreateStatusEndpoint(service)
		updatestatusEndpoint      = svc.MakeUpdateStatusEndpoint(service)
		deletestatusEndpoint      = svc.MakeDeleteStatusEndpoint(service)
		likestatusEndpoint        = svc.MakeLikeStatusEndpoint(service)
		unlikestatusEndpoint      = svc.MakeUnLikeStatusEndpoint(service)
		listlikestatusEndpoint    = svc.MakeListLikeStatusEndpoint(service)
		getstatusEndpoint         = svc.MakeGetStatusEndpoint(service)
		liststatusEndpoint        = svc.MakeListStatusEndpoint(service)
		listrecommendedEndpoint   = svc.MakeListRecommendedEndpoint(service)
		listusertimelineEndpoint  = svc.MakeListUserTimelineEndpoint(service)
		latestfollowingEndpoint   = svc.MakeLatestFollowingEndpoint(service)
		listrelationshipEndpoint  = svc.MakeListRelationshipEndpoint(service)
		followEndpoint            = svc.MakeFollowEndpoint(service)
		unfollowEndpoint          = svc.MakeUnFollowEndpoint(service)
		listmessageEndpoint       = svc.MakeListMessageEndpoint(service)
		readmessageEndpoint       = svc.MakeReadMessageEndpoint(service)
		getmessagesummaryEndpoint = svc.MakeGetMessageSummaryEndpoint(service)
		listcommentEndpoint       = svc.MakeListCommentEndpoint(service)
		createcommentEndpoint     = svc.MakeCreateCommentEndpoint(service)
		likecommentEndpoint       = svc.MakeLikeCommentEndpoint(service)
		unlikecommentEndpoint     = svc.MakeUnlikeCommentEndpoint(service)
		listblacklistEndpoint     = svc.MakeListBlacklistEndpoint(service)
		createblacklistEndpoint   = svc.MakeCreateBlacklistEndpoint(service)
		deleteblacklistEndpoint   = svc.MakeDeleteBlacklistEndpoint(service)
	)

	endpoints := svc.Endpoints{
		SignInEndpoint:            signinEndpoint,
		FindUserEndpoint:          finduserEndpoint,
		UpdateUserProfileEndpoint: updateuserprofileEndpoint,
		UpdateUserAvatarEndpoint:  updateuseravatarEndpoint,
		UpdateUserNameEndpoint:    updateusernameEndpoint,
		CreateStatusEndpoint:      createstatusEndpoint,
		UpdateStatusEndpoint:      updatestatusEndpoint,
		DeleteStatusEndpoint:      deletestatusEndpoint,
		LikeStatusEndpoint:        likestatusEndpoint,
		UnLikeStatusEndpoint:      unlikestatusEndpoint,
		ListLikeStatusEndpoint:    listlikestatusEndpoint,
		GetStatusEndpoint:         getstatusEndpoint,
		ListStatusEndpoint:        liststatusEndpoint,
		ListRecommendedEndpoint:   listrecommendedEndpoint,
		ListUserTimelineEndpoint:  listusertimelineEndpoint,
		LatestFollowingEndpoint:   latestfollowingEndpoint,
		ListRelationshipEndpoint:  listrelationshipEndpoint,
		FollowEndpoint:            followEndpoint,
		UnFollowEndpoint:          unfollowEndpoint,
		ListMessageEndpoint:       listmessageEndpoint,
		ReadMessageEndpoint:       readmessageEndpoint,
		GetMessageSummaryEndpoint: getmessagesummaryEndpoint,
		ListCommentEndpoint:       listcommentEndpoint,
		CreateCommentEndpoint:     createcommentEndpoint,
		LikeCommentEndpoint:       likecommentEndpoint,
		UnlikeCommentEndpoint:     unlikecommentEndpoint,
		ListBlacklistEndpoint:     listblacklistEndpoint,
		CreateBlacklistEndpoint:   createblacklistEndpoint,
		DeleteBlacklistEndpoint:   deleteblacklistEndpoint,
	}

	// Wrap selected Endpoints with middlewares. See handlers/middlewares.go
	endpoints = handlers.WrapEndpoints(endpoints)

	return endpoints
}

// Run starts a new http server, gRPC server, and a debug server with the
// passed config and logger
func Run(cfg svc.Config) {
	service := handlers.NewService()
	endpoints := NewEndpoints(service)

	if cfg.GenericHTTPResponseEncoder == nil {
		cfg.GenericHTTPResponseEncoder = svc.EncodeHTTPGenericResponse
	}

	// Mechanical domain.
	errc := make(chan error)

	// Interrupt handler.
	go handlers.InterruptHandler(errc)

	// Debug listener.
	go func() {
		log.Println("transport", "debug", "addr", cfg.DebugAddr)

		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

		errc <- http.ListenAndServe(cfg.DebugAddr, m)
	}()

	// HTTP transport.
	go func() {
		log.Println("transport", "HTTP", "addr", cfg.HTTPAddr)
		h := svc.MakeHTTPHandler(endpoints, cfg.GenericHTTPResponseEncoder)
		errc <- http.ListenAndServe(cfg.HTTPAddr, h)
	}()

	// gRPC transport.
	go func() {
		log.Println("transport", "gRPC", "addr", cfg.GRPCAddr)
		ln, err := net.Listen("tcp", cfg.GRPCAddr)
		if err != nil {
			errc <- err
			return
		}

		srv := svc.MakeGRPCServer(endpoints)
		s := grpc.NewServer()
		pb.RegisterSocialServer(s, srv)

		errc <- s.Serve(ln)
	}()

	// Run!
	log.Println("exit", <-errc)
}
