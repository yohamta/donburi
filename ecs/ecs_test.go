package ecs

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

func TestECS(t *testing.T) {
	world := donburi.NewWorld()
	ecs := NewECS(world)

	certainImage := ebiten.NewImage(1, 1)

	systems := []struct {
		system   interface{}
		image    *ebiten.Image
		priority int
	}{
		{&testSystem{}, nil, 0},
		{&testSystem{}, certainImage, 0},
		{&testDrawSystem{&testSystem{}}, nil, 1},
		{&testUpdateSystem{&testSystem{}}, nil, 1},
	}

	for _, s := range systems {
		opts := &SystemOpts{
			Priority: s.priority,
			Image:    s.image,
		}
		ecs.AddSystem(s.system, opts)
	}

	ecs.Update()

	updateTests := []struct {
		system               interface{}
		ExpectedUpdateCount  int
		ExpectedUpdatedIndex int
	}{
		{systems[0].system, 1, 1},
		{systems[1].system, 1, 2},
		{systems[3].system, 1, 0},
	}

	for idx, test := range updateTests {
		var sys *testSystem
		switch s := test.system.(type) {
		case *testSystem:
			sys = s
		case *testUpdateSystem:
			sys = s.system
		default:
			panic("invalid system")
		}

		if sys.UpdateCount != test.ExpectedUpdateCount {
			t.Errorf("test %d: expected update count %d, got %d", idx, test.ExpectedUpdateCount, sys.UpdateCount)
		}
		if sys.UpdatedIndex != test.ExpectedUpdatedIndex {
			t.Errorf("test %d: expected updated index %d, got %d", idx, test.ExpectedUpdatedIndex, sys.UpdatedIndex)
		}
	}

	defaultImage := ebiten.NewImage(1, 1)
	ecs.Draw(defaultImage)

	drawTests := []struct {
		system              interface{}
		ExpectedDrawCount   int
		ExpectedDrawedIndex int
		ExpectedImage       *ebiten.Image
	}{
		{systems[0].system, 1, 1, defaultImage},
		{systems[1].system, 1, 2, certainImage},
		{systems[2].system, 1, 0, defaultImage},
	}

	for idx, test := range drawTests {
		var sys *testSystem
		switch s := test.system.(type) {
		case *testSystem:
			sys = s
		case *testDrawSystem:
			sys = s.system
		default:
			panic("invalid system")
		}

		if sys.DrawCount != test.ExpectedDrawCount {
			t.Errorf("test %d: expected draw count %d, got %d", idx, test.ExpectedDrawCount, sys.DrawCount)
		}
		if sys.DrawedIndex != test.ExpectedDrawedIndex {
			t.Errorf("test %d: expected drawed index %d, got %d", idx, test.ExpectedDrawedIndex, sys.DrawedIndex)
		}
		if sys.DrawImage != test.ExpectedImage {
			t.Errorf("test %d: expected draw image %v, got %v", idx, test.ExpectedImage, sys.DrawImage)
		}
	}
}

var (
	testUpdatedIndex int
	testDrawedIndex  int
)

type testSystem struct {
	UpdatedIndex int
	DrawedIndex  int
	DrawImage    *ebiten.Image
	UpdateCount  int
	DrawCount    int
}

func (ts *testSystem) Update(ecs *ECS) {
	ts.UpdatedIndex = testUpdatedIndex
	ts.UpdateCount++

	testUpdatedIndex++
}

func (ts *testSystem) Draw(ecs *ECS, image *ebiten.Image) {
	ts.DrawedIndex = testDrawedIndex
	ts.DrawImage = image
	ts.DrawCount++

	testDrawedIndex++
}

type testUpdateSystem struct {
	system *testSystem
}

func (ts *testUpdateSystem) Update(ecs *ECS) {
	ts.system.Update(ecs)
}

type testDrawSystem struct {
	system *testSystem
}

func (ts *testDrawSystem) Draw(ecs *ECS, image *ebiten.Image) {
	ts.system.Draw(ecs, image)
}
