package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func unaryInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	s, ok := (*info).Server.(*Server)
	if !ok {
		return nil, fmt.Errorf("%T is not internalgrpc.Server", (*info).Server)
	}
	peer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, s.log.Error(ctx, "GRPC.Interceptor", fmt.Sprintf("can't get peer from incoming context %v", ctx))
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, s.log.Error(ctx, "GRPC.Interceptor", fmt.Sprintf("can't get metadata from incoming context %v", ctx))
	}

	sss := logger.GetCtxSession(ctx) // generate session uuid and fix start timestamp
	ctx = logger.PushCtxSession(ctx, sss)

	ua := md.Get("user-agent")
	if len(ua) == 0 {
		ua = []string{""}
	}
	s.log.Debug(ctx, "GRPC.Interceptor",
		fmt.Sprintf("New rpc-request to EventsService.[%s] from [%s] by [%s]",
			(*info).FullMethod, peer.Addr.String(), ua[0]),
		4,
	)
	sss.User, ok = interceptAuth(md)
	if !ok {
		err := status.Error(codes.Unauthenticated, "anonimus access prohibited")
		return nil, err
	}
	ctx = logger.PushCtxSession(ctx, sss)
	s.log.Debug(ctx, "GRPC.Auth", fmt.Sprintf("User [%s] auth success", sss.User), 2)

	if s.log.Level == "DEBUG" && s.log.Depth >= 5 {
		s.log.Debug(ctx, "GRPC.Interceptor", fmt.Sprintf("Request [%v]", req), 5)
	}

	respCode := 0
	resp, err = handler(ctx, req)
	stat, ok := status.FromError(err)
	if !ok {
		_ = s.log.Error(ctx, "GRPC.Interceptor", fmt.Sprintf("can't get status from error of handled req %v", err))
		respCode = -1
	} else {
		respCode = int(stat.Code())
	}

	interceptAfterResponse(s.log, ctx, respCode, resp)
	return resp, err
}

func interceptAuth(md metadata.MD) (string, bool) {
	user := md.Get("user")
	if len(user) == 0 {
		return "", false
	}
	return user[0], true
}

func interceptAfterResponse(log *logger.Logger, ctx context.Context, retCode int, response interface{}) {
	sss := logger.GetCtxSession(ctx)
	log.Debug(ctx, "GRPC.AfterResponse",
		fmt.Sprintf("Response status [%d] with latency %s", retCode, time.Since(sss.Start)),
		3,
	)
	log.Debug(ctx, "GRPC.AfterResponse", fmt.Sprintf("Response [%v]", response), 5)
}
