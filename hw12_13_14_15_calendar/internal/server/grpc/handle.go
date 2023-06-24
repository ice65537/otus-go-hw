package internalgrpc

import (
	"context"
	"fmt"

	eventsrv "github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/api"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ice65537/otus-go-hw/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) Hello(ctx context.Context, empty *emptypb.Empty) (hello *eventsrv.HelloMsg, err error) {
	sss := logger.GetCtxSession(ctx)
	hello = &eventsrv.HelloMsg{Text: fmt.Sprintf("Hello %s!", sss.User)}
	return
}

func evtFromReq(evtReq *eventsrv.Event) storage.Event {
	return storage.Event{
		ID:           evtReq.Id,
		Title:        evtReq.Title,
		StartDt:      evtReq.StartDt.AsTime(),
		StopDt:       evtReq.StopDt.AsTime(),
		Desc:         evtReq.Desc,
		Owner:        evtReq.Owner,
		NotifyBefore: int(evtReq.NotifyBefore),
	}
}

func (s *Server) New(ctx context.Context, evtReq *eventsrv.Event) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.edit(ctx, evtFromReq(evtReq), "gRPC.New")
}

func (s *Server) Reset(ctx context.Context, evtReq *eventsrv.Event) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.edit(ctx, evtFromReq(evtReq), "gRPC.Reset")
}

func (s *Server) Drop(ctx context.Context, guidReq *eventsrv.GUID) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, s.edit(ctx, storage.Event{ID: guidReq.Guid}, "gRPC.Drop")
}

func (s Server) edit(ctx context.Context, evt storage.Event, oper string) error {
	switch {
	case evt.ID == "" && oper == "gRPC.Drop":
		return status.Error(codes.InvalidArgument, "event.ID not found in request")
	case evt.ID != "" && oper == "gRPC.New":
		return status.Error(codes.InvalidArgument, "invalid path for existing event with ID, use /event/reset")
	case evt.Owner != "" && evt.Owner != logger.GetCtxSession(ctx).User:
		return status.Error(codes.PermissionDenied,
			fmt.Sprintf("%s, you can't create/update/drop event with Owner=[%s]", logger.GetCtxSession(ctx).User, evt.Owner))
	}
	evt.Owner = logger.GetCtxSession(ctx).User

	switch {
	case oper == "gRPC.New" || oper == "gRPC.Reset":
		if err := s.app.Upsert(ctx, evt); err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("%v", err))
		}
	case oper == "gRPC.Drop":
		if err := s.app.Drop(ctx, evt.ID); err != nil {
			return status.Error(codes.Internal, fmt.Sprintf("%v", err))
		}
	}
	return nil
}

func (s *Server) Get(ctx context.Context, evtSetReq *eventsrv.EventSetRequest) (evtSet *eventsrv.EventSet, err error) {
	events, err := s.app.Get(ctx, evtSetReq.T1.AsTime(), evtSetReq.T2.AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}
	evtSet = &eventsrv.EventSet{}
	evtSet.Events = make([]*eventsrv.Event, len(events))
	for i := range events {
		evtSet.Events[i] = &eventsrv.Event{
			Id:           events[i].ID,
			StartDt:      timestamppb.New(events[i].StartDt),
			StopDt:       timestamppb.New(events[i].StopDt),
			Desc:         events[i].Desc,
			Owner:        events[i].Owner,
			Title:        events[i].Title,
			NotifyBefore: int32(events[i].NotifyBefore),
		}
	}
	err = nil
	return
}
