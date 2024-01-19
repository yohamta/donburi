package ecs

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

func TestECS(t *testing.T) {
	world := donburi.NewWorld()
	ecs := NewECS(world)

	systems := []struct {
		layer  LayerID
		system *testSystem
	}{
		{1, &testSystem{}},
		{1, &testSystem{}},
		{0, &testSystem{}},
	}

	for _, sys := range systems {
		ecs.AddSystem(sys.system.Update)
		ecs.addRenderer(sys.layer, sys.system.Draw)
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

	ecs.Draw(ebiten.NewImage(1, 1))

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

func TestECSLayer(t *testing.T) {
	world := donburi.NewWorld()
	ecs := NewECS(world)

	var (
		layer0 LayerID = 0
		layer1 LayerID = 1
	)

	c1 := donburi.NewTag()

	ecs.Create(layer0, c1)
	ecs.Create(layer1, c1)

	systems := []struct {
		layer  LayerID
		system *testSystem
	}{
		{layer0, &testSystem{
			Query: NewQuery(layer0, filter.Contains(c1)),
		}},
		{layer1, &testSystem{
			Query: NewQuery(layer1, filter.Contains(c1)),
		}},
	}

	for _, sys := range systems {
		ecs.AddSystem(sys.system.Update)
		ecs.addRenderer(sys.layer, sys.system.Draw)
	}

	ecs.DrawLayer(0, ebiten.NewImage(1, 1))

	if systems[0].system.QueryCountDraw != 1 {
		t.Errorf("expected query count draw %d, got %d", 1, systems[0].system.QueryCountDraw)
	}
	if systems[1].system.QueryCountDraw != 0 {
		t.Errorf("expected query count draw %d, got %d", 0, systems[1].system.QueryCountDraw)
	}
}

func TestEmptyDefaultLayer(t *testing.T) {
	world := donburi.NewWorld()
	ecs := NewECS(world)

	TestLayer := LayerID(1)

	ecs.AddRenderer(TestLayer, func(ecs *ECS, image *ebiten.Image) {})

	ecs.Draw(ebiten.NewImage(1, 1))
}

var (
	testUpdatedIndex int
	testDrawedIndex  int
)

type testSystem struct {
	UpdatedIndex     int
	DrawedIndex      int
	DrawImage        *ebiten.Image
	UpdateCount      int
	DrawCount        int
	Query            *donburi.Query
	QueryCountUpdate int
	QueryCountDraw   int
}

func (ts *testSystem) Update(ecs *ECS) {
	ts.UpdatedIndex = testUpdatedIndex
	ts.UpdateCount++

	testUpdatedIndex++

	if ts.Query != nil {
		ts.QueryCountUpdate = ts.Query.Count(ecs.World)
	}
}

func (ts *testSystem) Draw(ecs *ECS, image *ebiten.Image) {
	ts.DrawedIndex = testDrawedIndex
	ts.DrawImage = image
	ts.DrawCount++

	testDrawedIndex++

	if ts.Query != nil {
		ts.QueryCountDraw = ts.Query.Count(ecs.World)
	}
}
