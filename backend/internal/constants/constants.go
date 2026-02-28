package constants

// Slot and time range constants for scheduling.
const (
	// SlotDurationSec is the length of one slot in seconds (15 min).
	SlotDurationSec = 15 * 60
	// MaxTimeUnix is used as exclusive end for "list all" range.
	MaxTimeUnix = int64(1) << 62
)

const (
	EventActionCreated = "created"
	EventActionDeleted = "deleted"
)

const (
	DurationSlotsMin = 1
	DurationSlotsMax = 2
)

const (
	DefaultDoctorID   = "dr1"
	DefaultDoctorName = "Dr. Smith"
)
