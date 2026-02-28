package service

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	schedulerv1 "scheduler/backend/gen/schedulerv1"
	"scheduler/backend/internal/broadcast"
	"scheduler/backend/internal/constants"
	"scheduler/backend/internal/repository"
	"scheduler/backend/internal/store"

	"github.com/google/uuid"
)

type Service struct {
	repo        repository.AppointmentRepository
	broadcaster *broadcast.Broadcaster
}

func New(repo repository.AppointmentRepository, broadcaster *broadcast.Broadcaster) *Service {
	return &Service{repo: repo, broadcaster: broadcaster}
}

func (s *Service) CreateAppointment(ctx context.Context, request *schedulerv1.CreateAppointmentRequest) (*schedulerv1.Appointment, error) {
	duration := request.DurationSlots
	if duration != constants.DurationSlotsMin && duration != constants.DurationSlotsMax {
		duration = constants.DurationSlotsMin
	}
	appointment := &schedulerv1.Appointment{
		Id:             uuid.New().String(),
		DoctorId:       request.DoctorId,
		DoctorName:     request.DoctorName,
		StartTimeUnix:  request.StartTimeUnix,
		DurationSlots:  duration,
		PatientName:    request.PatientName,
		PatientContact: request.PatientContact,
	}

	if err := s.repo.Create(ctx, appointment); err != nil {
		if err == store.ErrConflict {
			return nil, status.Error(codes.AlreadyExists, "One or more of the selected slots are already booked. Please choose a different time.")
		}
		log.Printf("service: create appointment: %v", err)
		return nil, status.Error(codes.Internal, "failed to create appointment")
	}
	s.broadcaster.Broadcast(&schedulerv1.AppointmentEvent{
		AppointmentId: appointment.Id,
		Action:        constants.EventActionCreated,
		Appointment:   appointment,
	})
	return appointment, nil
}

func (s *Service) ListAppointments(ctx context.Context) ([]*schedulerv1.Appointment, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		log.Printf("service: list appointments: %v", err)
		return nil, status.Error(codes.Internal, "failed to list appointments")
	}
	if list == nil {
		list = []*schedulerv1.Appointment{}
	}
	return list, nil
}

func (s *Service) DeleteAppointment(ctx context.Context, id string) error {
	existed, err := s.repo.Delete(ctx, id)
	if err != nil {
		log.Printf("service: delete appointment %s: %v", id, err)
		return status.Error(codes.Internal, "failed to delete appointment")
	}
	if existed {
		s.broadcaster.Broadcast(&schedulerv1.AppointmentEvent{
			AppointmentId: id,
			Action:        constants.EventActionDeleted,
		})
	}
	return nil
}
