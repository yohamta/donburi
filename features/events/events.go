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
		eventName     string
		eventBus      donburi.ComponentType
		eventBusQuery *query.Query
	}

	// Subscriber is a function that handles an event.
	Subscriber[T any] func(w donburi.World, event T)
)

type (
	eventBusData[T any] struct {
		subscribers []Subscriber[T]
		queue       []T
	}
	eventTypeData struct {
		eventName string
		process   func(w donburi.World)
	}
)

var (
	eventType  = donburi.NewComponentType[eventTypeData]()
	eventQuery = query.NewQuery(filter.Contains(eventType))
)

// ProcessAllEvents processes all events.
func ProcessAllEvents(w donburi.World) {
	eventQuery.EachEntity(w, func(e *donburi.Entry) {
		eventType := getEventType(e)
		eventType.process(w)
	})
}

// NewEventType creates a new event type.
func NewEventType[T any]() *EventType[T] {
	eventBus := donburi.NewComponentType[eventBusData[T]]()
	var t T
	e := &EventType[T]{
		eventName:     reflect.TypeOf(t).Name(),
		eventBus:      eventBus,
		eventBusQuery: query.NewQuery(filter.Contains(eventBus)),
	}
	return e
}

// RegisterHandler registers a subscriber for the event.
func (e *EventType[T]) Subscribe(w donburi.World, subscriber Subscriber[T]) {
	if e.eventBusQuery.Count(w) == 0 {
		entity := w.Entry(w.Create(e.eventBus, eventType))
		donburi.Set(entity, e.eventBus, newEventBusData[T]())
		donburi.SetValue(entity, eventType, eventTypeData{
			eventName: e.eventName,
			process: func(w donburi.World) {
				e.ProcessEvents(w)
			},
		})
	}

	eventBus := e.mustFindEventBus(w)
	eventBus.subscribers = append(eventBus.subscribers, subscriber)
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
			for _, s := range eventBus.subscribers {
				if Debug {
					fmt.Printf("%T -> %T\n", event, s)
				}

				s(w, event)
			}
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

func getEventType(e *donburi.Entry) *eventTypeData {
	return donburi.Get[eventTypeData](e, eventType)
}

func newEventBusData[T any]() *eventBusData[T] {
	return &eventBusData[T]{}
}
