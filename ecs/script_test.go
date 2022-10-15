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

	entityA := world.Create(testEntityA)
	entityB := world.Create(testEntityB)

	queryA := query.NewQuery(filter.Contains(testEntityA))
	queryB := query.NewQuery(filter.Contains(testEntityB))

	scriptA := &testScript{}
	scriptB := &testScript{}

	ecs.AddUpdateScript(scriptA, queryA)
	ecs.AddUpdateScript(scriptB, queryB)

	ecs.AddDrawScript(0, scriptA, queryA)
	ecs.AddDrawScript(0, scriptB, queryB)

	ecs.Update()

	ecs.Draw(0, ebiten.NewImage(1, 1))

	tests := []struct {
		script          *testScript
		entity          donburi.Entity
		expectedUpdated int
		expectedDrawed  int
	}{
		{
			script:          scriptA,
			entity:          entityA,
			expectedUpdated: 1,
			expectedDrawed:  1,
		},
		{
			script:          scriptB,
			entity:          entityB,
			expectedUpdated: 1,
			expectedDrawed:  1,
		},
	}

	for idx, test := range tests {
		scr := test.script
		if scr.UpdateEntry.Entity() != test.entity {
			t.Errorf("test %d: expected update entry entity %v, got %v", idx, test.entity, scr.UpdateEntry.Entity())
		}
		if scr.UpdatedCount != test.expectedUpdated {
			t.Errorf("test %d: expected updated count %d, got %d", idx, test.expectedUpdated, scr.UpdatedCount)
		}
		if scr.DrawEntry.Entity() != test.entity {
			t.Errorf("test %d: expected draw entry entity %v, got %v", idx, test.entity, scr.DrawEntry.Entity())
		}
		if scr.DrawCount != test.expectedDrawed {
			t.Errorf("test %d: expected draw count %d, got %d", idx, test.expectedDrawed, scr.DrawCount)
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
