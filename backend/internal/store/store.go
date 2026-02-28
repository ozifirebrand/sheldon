package store

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	schedulerv1 "scheduler/backend/gen/schedulerv1"
	"scheduler/backend/internal/constants"

	_ "modernc.org/sqlite"
)

var ErrConflict = errors.New("scheduling conflict: one or more slots already booked")

const schema = `
CREATE TABLE IF NOT EXISTS appointments (
  id TEXT PRIMARY KEY,
  doctor_id TEXT NOT NULL,
  doctor_name TEXT NOT NULL,
  start_time_unix INTEGER NOT NULL,
  duration_slots INTEGER NOT NULL,
  patient_name TEXT NOT NULL,
  patient_contact TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_appointments_doctor_slot ON appointments(doctor_id, start_time_unix);
`

type Store struct {
	db *sql.DB
	mu sync.Mutex
}

func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Create(ctx context.Context, appointment *schedulerv1.Appointment) error {
	appointmentDuration := appointment.DurationSlots
	if appointmentDuration != constants.DurationSlotsMin && appointmentDuration != constants.DurationSlotsMax {
		appointmentDuration = constants.DurationSlotsMin
	}
	appointment.DurationSlots = appointmentDuration

	s.mu.Lock()
	defer s.mu.Unlock()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := int32(0); i < appointment.DurationSlots; i++ {
		slotStart := appointment.StartTimeUnix + int64(i)*constants.SlotDurationSec
		var existing string
		err := tx.QueryRowContext(ctx,
			`SELECT id FROM appointments WHERE doctor_id = ? AND start_time_unix = ?`,
			appointment.DoctorId, slotStart).Scan(&existing)
		if err == nil {
			return ErrConflict
		}
		if err != sql.ErrNoRows {
			return err
		}
	}
	endingUnix := appointment.StartTimeUnix + int64(appointment.DurationSlots)*constants.SlotDurationSec
	var count int
	err = tx.QueryRowContext(ctx,
		`SELECT COUNT(1) FROM appointments WHERE doctor_id = ? AND
		 start_time_unix + duration_slots * ? > ? AND start_time_unix < ?`,
		appointment.DoctorId, constants.SlotDurationSec, appointment.StartTimeUnix, endingUnix).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrConflict
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO appointments (id, doctor_id, doctor_name, start_time_unix, duration_slots, patient_name, patient_contact)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		appointment.Id, appointment.DoctorId, appointment.DoctorName, appointment.StartTimeUnix, appointment.DurationSlots, appointment.PatientName, appointment.PatientContact)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) List(ctx context.Context) ([]*schedulerv1.Appointment, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, doctor_id, doctor_name, start_time_unix, duration_slots, patient_name, patient_contact
		 FROM appointments ORDER BY start_time_unix`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*schedulerv1.Appointment
	for rows.Next() {
		var appointment schedulerv1.Appointment
		if err := rows.Scan(&appointment.Id, &appointment.DoctorId, &appointment.DoctorName, &appointment.StartTimeUnix, &appointment.DurationSlots, &appointment.PatientName, &appointment.PatientContact); err != nil {
			return nil, err
		}
		appointmentCopy := appointment
		out = append(out, &appointmentCopy)
	}
	return out, rows.Err()
}

func (s *Store) Delete(ctx context.Context, id string) (existed bool, err error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM appointments WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (s *Store) Get(ctx context.Context, id string) (*schedulerv1.Appointment, error) {
	var appointment schedulerv1.Appointment
	err := s.db.QueryRowContext(ctx,
		`SELECT id, doctor_id, doctor_name, start_time_unix, duration_slots, patient_name, patient_contact
		 FROM appointments WHERE id = ?`, id).Scan(
		&appointment.Id, &appointment.DoctorId, &appointment.DoctorName, &appointment.StartTimeUnix, &appointment.DurationSlots, &appointment.PatientName, &appointment.PatientContact)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}
