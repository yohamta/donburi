package ecs

import (
	"reflect"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

// Debug enables debug logging.
var Debug = false

// Command represents kind of command.
// It is used to register handler and publish commands.
type (
	Command[T any] struct {
		cmdName  string
		busType  *donburi.ComponentType[T]
		busQuery *donburi.Query
	}

	// Handler is a function that handles an event.
	Handler[T any] func(ecs *ECS, command T)
)

type (
	commandBusData[T any] struct {
		handler Handler[T]
		queue   []T
	}
	commandTypeData struct {
		commandName string
		process     func(ecs *ECS)
	}
)

var (
	commandType  = donburi.NewComponentType[commandTypeData]()
	commandQuery = donburi.NewQuery(filter.Contains(commandType))
)

// ProcessAllCommands processes all events.
func ProcessAllCommands(ecs *ECS) {
	commandQuery.Each(ecs.World, func(e *donburi.Entry) {
		cmdType := getCommandType(e)
		cmdType.process(ecs)
	})
}

// NewCommand creates a new command.
func NewCommand[T any]() *Command[T] {
	busType := donburi.NewComponentType[T]()
	var t T
	e := &Command[T]{
		cmdName:  reflect.TypeOf(t).Name(),
		busType:  busType,
		busQuery: donburi.NewQuery(filter.Contains(busType)),
	}
	return e
}

// Register registers a handler for the command.
func (e *Command[T]) Register(ecs *ECS, handler Handler[T]) {
	if e.busQuery.Count(ecs.World) == 0 {
		entity := ecs.World.Entry(ecs.World.Create(e.busType, commandType))
		donburi.Set(entity, e.busType, newBusData[T]())
		donburi.SetValue(entity, commandType, commandTypeData{
			commandName: e.cmdName,
			process: func(ecs *ECS) {
				e.ProcessCommands(ecs)
			},
		})
	}

	bus := e.mustFindBus(ecs.World)
	bus.handler = handler
}

// Dispatch dipatches a command.
func (e *Command[T]) Dispatch(ecs *ECS, event T) {
	bus := e.mustFindBus(ecs.World)
	bus.handler(ecs, event)
	bus.queue = append(bus.queue, event)
}

// ProcessCommands processes commands.
func (e *Command[T]) ProcessCommands(ecs *ECS) {
	bus := e.mustFindBus(ecs.World)
	for len(bus.queue) > 0 {
		queue := bus.queue
		bus.queue = nil
		for _, cmd := range queue {
			bus.handler(ecs, cmd)
		}
	}
}

func (e *Command[T]) mustFindBus(w donburi.World) *commandBusData[T] {
	bus, ok := e.busQuery.First(w)
	if !ok {
		panic("command bus not found")
	}
	return donburi.Get[commandBusData[T]](bus, e.busType)
}

func getCommandType(e *donburi.Entry) *commandTypeData {
	return donburi.Get[commandTypeData](e, commandType)
}

func newBusData[T any]() *commandBusData[T] {
	return &commandBusData[T]{}
}
