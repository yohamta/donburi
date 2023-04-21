package events_test

import (
	"testing"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type EnemyKilled struct {
	EnemyID int
}

var EnemyKilledEvent = events.NewEventType[*EnemyKilled]()

func TestEvents(t *testing.T) {

	w := donburi.NewWorld()

	EnemyKilledEvent.Subscribe(w, HandleEnemyKilled)
	EnemyKilledEvent.Publish(w, &EnemyKilled{EnemyID: 1})

	events.ProcessAllEvents(w)

	ev := lastReceivedEvent

	if ev == nil {
		t.Errorf("event should be received")
	}

	if ev.EnemyID != 1 {
		t.Errorf("event should have value 1")
	}

	lastReceivedEvent = nil
	events.ProcessAllEvents(w)

	if lastReceivedEvent != nil {
		t.Errorf("event should not be received")
	}
}

func TestUnsubscribe(t *testing.T) {

	w := donburi.NewWorld()

	EnemyKilledEvent.Subscribe(w, HandleEnemyKilled)
	EnemyKilledEvent.Publish(w, &EnemyKilled{EnemyID: 1})

	events.ProcessAllEvents(w)

	ev := lastReceivedEvent

	if ev == nil {
		t.Errorf("event should be received")
	}

	if ev.EnemyID != 1 {
		t.Errorf("event should have value 1")
	}

	lastReceivedEvent = nil
	EnemyKilledEvent.Unsubscribe(w, HandleEnemyKilled)
	EnemyKilledEvent.Publish(w, &EnemyKilled{EnemyID: 1})

	events.ProcessAllEvents(w)

	if lastReceivedEvent != nil {
		t.Errorf("event should not be received")
	}
}

func TestEventsWithNoSubscribers(t *testing.T) {

	w := donburi.NewWorld()

	EnemyKilledEvent.Publish(w, &EnemyKilled{EnemyID: 1})

	events.ProcessAllEvents(w)

	ev := lastReceivedEvent

	if ev != nil {
		t.Errorf("event should not be received")
	}
}

var (
	lastReceivedEvent *EnemyKilled = nil
)

func HandleEnemyKilled(w donburi.World, event *EnemyKilled) {
	lastReceivedEvent = event
}
