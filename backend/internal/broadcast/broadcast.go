package broadcast

import (
	"sync"

	schedulerv1 "scheduler/backend/gen/schedulerv1"
)

// Broadcaster sends appointment events to all subscribed streams.
type Broadcaster struct {
	mu      sync.Mutex
	streams []chan *schedulerv1.AppointmentEvent
}

// New creates a new Broadcaster.
func New() *Broadcaster {
	return &Broadcaster{streams: make([]chan *schedulerv1.AppointmentEvent, 0)}
}

// Subscribe adds a channel that will receive all events until Unsubscribe is called.
func (b *Broadcaster) Subscribe() (ch <-chan *schedulerv1.AppointmentEvent, unsubscribe func()) {
	c := make(chan *schedulerv1.AppointmentEvent, 16)
	b.mu.Lock()
	b.streams = append(b.streams, c)
	idx := len(b.streams) - 1
	b.mu.Unlock()
	unsubscribe = func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if idx < len(b.streams) {
			b.streams[idx] = b.streams[len(b.streams)-1]
			b.streams = b.streams[:len(b.streams)-1]
		}
		close(c)
	}
	return c, unsubscribe
}

// Broadcast sends the event to all current subscribers. Non-blocking; drops if channel full.
func (b *Broadcaster) Broadcast(ev *schedulerv1.AppointmentEvent) {
	b.mu.Lock()
	streams := make([]chan *schedulerv1.AppointmentEvent, len(b.streams))
	copy(streams, b.streams)
	b.mu.Unlock()
	for _, ch := range streams {
		select {
		case ch <- ev:
		default:
		}
	}
}
