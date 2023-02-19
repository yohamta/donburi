package ecs_test

import (
	"testing"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type DoSomthing struct {
	SomeParam int
}

var DoSomething = ecs.NewCommand[*DoSomthing]()

func TestEvents(t *testing.T) {

	e := ecs.NewECS(donburi.NewWorld())

	DoSomething.Register(e, HandleDoSomething)
	DoSomething.Dispatch(e, &DoSomthing{SomeParam: 1})

	ecs.ProcessAllCommands(e)

	ev := lastReceivedCommand

	if ev == nil {
		t.Errorf("cmd should be received")
	} else if ev.SomeParam != 1 {
		t.Errorf("cmd should have value 1")
	}

	lastReceivedCommand = nil
	ecs.ProcessAllCommands(e)

	if lastReceivedCommand != nil {
		t.Errorf("event should not be received")
	}
}

var (
	lastReceivedCommand *DoSomthing = nil
)

func HandleDoSomething(e *ecs.ECS, cmd *DoSomthing) {
	lastReceivedCommand = cmd
}
