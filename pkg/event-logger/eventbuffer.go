package elogger

import (
	"container/ring"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// The EventBuffer interface is a basic interface to interact with a buffer
// for storing events
type EventBuffer interface {
	Add(*corev1.Event)
	Do(f func(*corev1.Event))
	Capacity() int
	Size() int
}

// The RingEventBuffer is a simple deduplicating buffer to store events in a ring,
// the ring structure means old events will be overwritten by new events.
type RingEventBuffer struct {
	r *ring.Ring
	s map[types.UID]bool
}

// NewRingEventBuffer creates a new event buffer of size `bufferSize`
func NewRingEventBuffer(bufferSize int) *RingEventBuffer {
	rv := RingEventBuffer{
		r: ring.New(bufferSize),
		s: make(map[types.UID]bool),
	}

	return &rv
}

// Add add's an event to the buffer
func (b *RingEventBuffer) Add(e *corev1.Event) {
	if _, exists := b.s[e.UID]; exists {
		return
	}

	if b.r.Value != nil {
		uid := b.r.Value.(*corev1.Event).UID
		delete(b.s, uid)
	}

	b.s[e.UID] = true

	b.r.Value = e
	b.r = b.r.Next()
}

// Do performs a function on all events in the buffer
func (b *RingEventBuffer) Do(f func(*corev1.Event)) {
	b.r.Do(func(v any) {
		if v == nil {
			return
		}
		f(v.(*corev1.Event))
	})
}

// Capacity returns the max capacity of the buffer
func (b *RingEventBuffer) Capacity() int {
	return b.r.Len()
}

// Size returns the number of events currently in the buffer
func (b *RingEventBuffer) Size() int {
	return len(b.s)
}
