package events

import (
	"fmt"
	"reflect"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

// Debug enables debug logging.
var Debug = false

// Event represents kind of event.
// It is used to subscribe and publish events.
type (
	EventType[T any] struct {
		Subscriber Subscriber[T]

		eventName     string
		eventId       eventId
		eventBus      *donburi.ComponentType
		eventBusQuery *query.Query
	}

	// Subscriber is a function that handles an event.
	Subscriber[T any] func(w donburi.World, event T)
)

type (
	eventId             int
	eventBusData[T any] struct {
		subscribers Subscriber[T]
		queue       []T
	}
	eventType struct {
		process func(w donburi.World)
	}
)

var (
	registeredEvents         = []*eventType{}
	nextId           eventId = 0
)

// NewEventType creates a new event type.
func NewEventType[T any](subscriber Subscriber[T]) *EventType[T] {
	eventBus := donburi.NewComponentType[eventBusData[T]]()
	var t T
	e := &EventType[T]{
		Subscriber:    subscriber,
		eventName:     reflect.TypeOf(t).Name(),
		eventId:       nextId,
		eventBus:      eventBus,
		eventBusQuery: query.NewQuery(filter.Contains(eventBus)),
	}
	nextId++
	registeredEvents = append(registeredEvents, &eventType{
		process: func(w donburi.World) {
			e.ProcessEvents(w)
		},
	})
	donburi.RegisterInitializer(
		func(w donburi.World) {
			entity := w.Entry(w.Create(eventBus))
			donburi.Set(entity, eventBus, newEventBusData[T]())
		},
	)
	return e
}

// ProcessAllEvents processes all events.
func ProcessAllEvents(w donburi.World) {
	for _, e := range registeredEvents {
		e.process(w)
	}
}

// Publish publishes an event.
func (e *EventType[T]) Publish(w donburi.World, event T) {
	eventBus := e.mustFindEventBus(w)
	if Debug {
		fmt.Printf("Publishing %T\n", event)
	}
	eventBus.queue = append(eventBus.queue, event)
}

// ProcessEvents processes events.
func (e *EventType[T]) ProcessEvents(w donburi.World) {
	eventBus := e.mustFindEventBus(w)
	// The outer loop is needed, because events can trigger more events.
	for len(eventBus.queue) > 0 {
		queue := eventBus.queue
		eventBus.queue = nil
		for _, event := range queue {
			if Debug {
				fmt.Printf("%T -> %T\n", event, e.Subscriber)
			}

			e.Subscriber(w, event)
		}
	}
}

func (e *EventType[T]) mustFindEventBus(w donburi.World) *eventBusData[T] {
	eventBus, ok := e.eventBusQuery.FirstEntity(w)
	if !ok {
		panic("event bus not found")
	}
	return donburi.Get[eventBusData[T]](eventBus, e.eventBus)
}

func newEventBusData[T any]() *eventBusData[T] {
	return &eventBusData[T]{}
}
