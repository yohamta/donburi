package ecs

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

var (
	testEntityA = donburi.NewTag()
	testEntityB = donburi.NewTag()
)

func TestScriptSystem(t *testing.T) {
	world := donburi.NewWorld()
	ecs := NewECS(world)

	certainImage := ebiten.NewImage(1, 1)

	entityA := world.Create(testEntityA)
	entityB := world.Create(testEntityB)

	scriptA := &testScript{}
	scriptB := &testScript{}

	ecs.AddScript(*query.NewQuery(
		filter.Contains(testEntityA),
	), scriptA, &ScriptOpts{
		Image: certainImage,
	})

	ecs.AddScript(*query.NewQuery(
		filter.Contains(testEntityB),
	), scriptB, nil)

	ecs.Update()

	defaultImage := ebiten.NewImage(1, 1)
	ecs.Draw(defaultImage)

	tests := []struct {
		script          *testScript
		entity          donburi.Entity
		expectedUpdated int
		expectedDrawed  int
		expectedImage   *ebiten.Image
	}{
		{
			script:          scriptA,
			entity:          entityA,
			expectedUpdated: 1,
			expectedDrawed:  1,
			expectedImage:   certainImage,
		},
		{
			script:          scriptB,
			entity:          entityB,
			expectedUpdated: 1,
			expectedDrawed:  1,
			expectedImage:   defaultImage,
		},
	}

	for idx, test := range tests {
		if test.script.UpdateEntry.Entity() != test.entity {
			t.Errorf("test %d: expected update entry entity %v, got %v", idx, test.entity, test.script.UpdateEntry.Entity())
		}
		if test.script.UpdatedCount != test.expectedUpdated {
			t.Errorf("test %d: expected updated count %d, got %d", idx, test.expectedUpdated, test.script.UpdatedCount)
		}
		if test.script.DrawEntry.Entity() != test.entity {
			t.Errorf("test %d: expected draw entry entity %v, got %v", idx, test.entity, test.script.DrawEntry.Entity())
		}
		if test.script.DrawCount != test.expectedDrawed {
			t.Errorf("test %d: expected draw count %d, got %d", idx, test.expectedDrawed, test.script.DrawCount)
		}
	}
}

type testScript struct {
	UpdatedCount int
	DrawCount    int
	UpdateEntry  *donburi.Entry
	DrawEntry    *donburi.Entry
}

func (ts *testScript) Update(entry *donburi.Entry) {
	ts.UpdateEntry = entry
	ts.UpdatedCount++
}

func (ts *testScript) Draw(entry *donburi.Entry, screen *ebiten.Image) {
	ts.DrawEntry = entry
	ts.DrawCount++
}