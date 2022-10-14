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
		system   *testSystem
		image    *ebiten.Image
		priority int
	}{
		{&testSystem{}, nil, 0},
		{&testSystem{}, certainImage, 0},
		{&testSystem{}, nil, 1},
	}

	for _, sys := range systems {
		ecs.AddSystem(sys.system, &SystemOpts{
			Priority: sys.priority,
			Image:    sys.image,
		})
	}

	ecs.Update()

	updateTests := []struct {
		system               *testSystem
		ExpectedUpdateCount  int
		ExpectedUpdatedIndex int
	}{
		{systems[0].system, 1, 1},
		{systems[1].system, 1, 2},
		{systems[2].system, 1, 0},
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

	defaultImage := ebiten.NewImage(1, 1)
	ecs.Draw(defaultImage)

	drawTests := []struct {
		system              *testSystem
		ExpectedDrawCount   int
		ExpectedDrawedIndex int
		ExpectedImage       *ebiten.Image
	}{
		{systems[0].system, 1, 1, defaultImage},
		{systems[1].system, 1, 2, certainImage},
		{systems[2].system, 1, 0, defaultImage},
	}

	for idx, test := range drawTests {
		sys := test.system
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
