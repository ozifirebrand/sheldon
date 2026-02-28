package constants

// Slot and time range constants for scheduling.
const (
	// SlotDurationSec is the length of one slot in seconds (15 min).
	SlotDurationSec = 15 * 60
	// MaxTimeUnix is used as exclusive end for "list all" range.
	MaxTimeUnix = int64(1) << 62
)

// Event action strings for appointment lifecycle.
const (
	EventActionCreated = "created"
	EventActionDeleted = "deleted"
)

// Allowed duration in slots (1 = 15 min, 2 = 30 min).
const (
	DurationSlotsMin = 1
	DurationSlotsMax = 2
)

// Default doctor for single-doctor assessment (e.g. GET /api/doctors).
const (
	DefaultDoctorID   = "dr1"
	DefaultDoctorName = "Dr. Smith"
)
