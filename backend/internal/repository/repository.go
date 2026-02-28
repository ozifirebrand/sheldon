package repository

import (
	"context"

	schedulerv1 "scheduler/backend/gen/schedulerv1"
)

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *schedulerv1.Appointment) error
	List(ctx context.Context) ([]*schedulerv1.Appointment, error)
	Get(ctx context.Context, id string) (*schedulerv1.Appointment, error)
	Delete(ctx context.Context, id string) (existed bool, err error)
}
