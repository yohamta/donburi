package events_test

import (
	"testing"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type anEvent struct {
	Value int
}

var anEventType = events.NewEventType(AnEventHandler)

func TestEvents(t *testing.T) {
	w := donburi.NewWorld()

	anEventType.Publish(w, &anEvent{Value: 1})

	events.ProcessAllEvents(w)

	ev := lastReceivedEvent

	if ev == nil {
		t.Errorf("event should be received")
	}

	if ev.Value != 1 {
		t.Errorf("event should have value 1")
	}
}

var (
	lastReceivedEvent *anEvent = nil
)

func AnEventHandler(w donburi.World, event *anEvent) {
	lastReceivedEvent = event
}
