package events

import (
	"fmt"
	"reflect"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// Debug enables debug logging.
var Debug = false

// Event represents kind of event.
// It is used to subscribe and publish events.
type (
	EventType[T any] struct {
		eventName     string
		eventBus      *donburi.ComponentType[T]
		eventBusQuery *donburi.Query
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
	eventQuery = donburi.NewQuery(filter.Contains(eventType))
)

// ProcessAllEvents processes all events.
func ProcessAllEvents(w donburi.World) {
	eventQuery.Each(w, func(e *donburi.Entry) {
		eventType := getEventType(e)
		eventType.process(w)
	})
}

// NewEventType creates a new event type.
func NewEventType[T any]() *EventType[T] {
	eventBus := donburi.NewComponentType[T]()
	var t T
	e := &EventType[T]{
		eventName:     reflect.TypeOf(t).Name(),
		eventBus:      eventBus,
		eventBusQuery: donburi.NewQuery(filter.Contains(eventBus)),
	}
	return e
}

// Subscribe registers a subscriber for the event.
func (e *EventType[T]) Subscribe(w donburi.World, subscriber Subscriber[T]) {
	eventBus := e.mustFindEventBus(w)
	eventBus.subscribers = append(eventBus.subscribers, subscriber)
}

// Find the index of a subscriber
func (e *EventType[T]) findSubscriber(w donburi.World, subscriber Subscriber[T]) (int, bool) {
	if e.eventBusQuery.Count(w) != 0 {
		eventBus := e.mustFindEventBus(w)
		for i, s := range eventBus.subscribers {
			if reflect.ValueOf(s).Pointer() == reflect.ValueOf(subscriber).Pointer() {
				return i, true
			}
		}
	}
	return 0, false
}

// Unsubscribe removes a subscriber for the event.
func (e *EventType[T]) Unsubscribe(w donburi.World, subscriber Subscriber[T]) {
	index, found := e.findSubscriber(w, subscriber)
	if found {
		eventBus := e.mustFindEventBus(w)
		eventBus.subscribers[index] = eventBus.subscribers[len(eventBus.subscribers)-1]
		eventBus.subscribers = eventBus.subscribers[:len(eventBus.subscribers)-1]
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
	eventBus, ok := e.eventBusQuery.First(w)
	if !ok {
		eventBus = w.Entry(w.Create(e.eventBus, eventType))
		donburi.Set(eventBus, e.eventBus, newEventBusData[T]())
		donburi.SetValue(eventBus, eventType, eventTypeData{
			eventName: e.eventName,
			process: func(w donburi.World) {
				e.ProcessEvents(w)
			},
		})
	}
	return donburi.Get[eventBusData[T]](eventBus, e.eventBus)
}

func getEventType(e *donburi.Entry) *eventTypeData {
	return donburi.Get[eventTypeData](e, eventType)
}

func newEventBusData[T any]() *eventBusData[T] {
	return &eventBusData[T]{}
}
