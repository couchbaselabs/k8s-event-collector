package elogger

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/rand"
)

func createEvent() corev1.Event {
	k := corev1.Event{
		ObjectMeta: v1.ObjectMeta{
			UID:             types.UID(rand.String(36)),
			ResourceVersion: "1",
		},
	}
	return k
}

func TestDeduplication(t *testing.T) {
	e := createEvent()

	b := NewRingEventBuffer(4)
	b.Add(&e)
	b.Add(&e)
	b.Add(&e)

	eventsInBuffer := 0

	b.Do(func(e *corev1.Event) {
		eventsInBuffer++
	})
	if eventsInBuffer != 1 {
		t.Errorf("Events should have been de-duplicated, only expecting one event")
	}
}

func TestEventTracking(t *testing.T) {
	bufferSize := 4
	b := NewRingEventBuffer(bufferSize)

	for i := 0; i < 10; i++ {
		e := createEvent()
		b.Add(&e)
	}

	if b.r.Len() != bufferSize {
		t.Errorf("The buffer ring should be of size: %v", bufferSize)
	}

	if len(b.s) != bufferSize {
		t.Errorf("The buffer set should be of size: %v", bufferSize)
	}
}
