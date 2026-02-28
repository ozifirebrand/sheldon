package server

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	schedulerv1 "scheduler/backend/gen/schedulerv1"
	"scheduler/backend/internal/broadcast"
	"scheduler/backend/internal/service"
)

type Server struct {
	schedulerv1.UnimplementedSchedulerServiceServer
	schedulerService *service.Service
	broadcaster      *broadcast.Broadcaster
}

func New(schedulerService *service.Service, broadcaster *broadcast.Broadcaster) *Server {
	return &Server{schedulerService: schedulerService, broadcaster: broadcaster}
}

func (s *Server) CreateAppointment(ctx context.Context, req *schedulerv1.CreateAppointmentRequest) (*schedulerv1.Appointment, error) {
	return s.schedulerService.CreateAppointment(ctx, req)
}

func (s *Server) ListAppointments(ctx context.Context, req *schedulerv1.ListAppointmentsRequest) (*schedulerv1.ListAppointmentsResponse, error) {
	list, err := s.schedulerService.ListAppointments(ctx)
	if err != nil {
		return nil, err
	}
	return &schedulerv1.ListAppointmentsResponse{Appointments: list}, nil
}

func (s *Server) DeleteAppointment(ctx context.Context, req *schedulerv1.DeleteAppointmentRequest) (*schedulerv1.DeleteAppointmentResponse, error) {
	if err := s.schedulerService.DeleteAppointment(ctx, req.Id); err != nil {
		return nil, err
	}
	return &schedulerv1.DeleteAppointmentResponse{}, nil
}

func (s *Server) SubscribeAppointmentUpdates(req *schedulerv1.SubscribeAppointmentUpdatesRequest, stream grpc.ServerStreamingServer[schedulerv1.AppointmentEvent]) error {
	ch, unsubscribe := s.broadcaster.Subscribe()
	defer unsubscribe()
	for ev := range ch {
		if err := stream.Send(ev); err != nil {
			log.Printf("grpc: stream send: %v", err)
			return status.Error(codes.Internal, "stream send failed")
		}
	}
	return nil
}
