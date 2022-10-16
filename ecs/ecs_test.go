package ecs

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

func TestECS(t *testing.T) {
	world := donburi.NewWorld()
	ecs := NewECS(world)

	systems := []struct {
		layer  DrawLayer
		system *testSystem
	}{
		{1, &testSystem{}},
		{1, &testSystem{}},
		{0, &testSystem{}},
	}

	for _, sys := range systems {
		ecs.AddSystem(System{
			Update:    sys.system.Update,
			DrawLayer: sys.layer,
			Draw:      sys.system.Draw,
		})
	}

	ecs.Update()

	updateTests := []struct {
		system               *testSystem
		ExpectedUpdateCount  int
		ExpectedUpdatedIndex int
	}{
		{systems[0].system, 1, 0},
		{systems[1].system, 1, 1},
		{systems[2].system, 1, 2},
	}

	for idx, test := range updateTests {
		sys := test.system
		if sys.UpdateCount != test.ExpectedUpdateCount {
			t.Errorf("test %d: expected update count %d, got %d", idx, test.ExpectedUpdateCount, sys.UpdateCount)
		}
		if sys.UpdatedIndex != test.ExpectedUpdatedIndex {
			t.Errorf("test %d: expected updated index %d, got %d", idx, test.ExpectedUpdatedIndex, sys.UpdatedIndex)
		}
	}

	ecs.Draw(0, ebiten.NewImage(1, 1))
	ecs.Draw(1, ebiten.NewImage(1, 1))

	drawTests := []struct {
		system              *testSystem
		ExpectedDrawCount   int
		ExpectedDrawedIndex int
	}{
		{systems[0].system, 1, 1},
		{systems[1].system, 1, 2},
		{systems[2].system, 1, 0},
	}

	for idx, test := range drawTests {
		sys := test.system
		if sys.DrawCount != test.ExpectedDrawCount {
			t.Errorf("test %d: expected draw count %d, got %d", idx, test.ExpectedDrawCount, sys.DrawCount)
		}
		if sys.DrawedIndex != test.ExpectedDrawedIndex {
			t.Errorf("test %d: expected drawed index %d, got %d", idx, test.ExpectedDrawedIndex, sys.DrawedIndex)
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
