package donburi_test

import (
	"testing"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var (
	queryTagA = donburi.NewTag()
	queryTagB = donburi.NewTag()
	queryTagC = donburi.NewTag()
)

func TestQuery(t *testing.T) {
	world := donburi.NewWorld()
	world.Create(queryTagA)
	world.Create(queryTagC)
	world.Create(queryTagA, queryTagB)

	query := donburi.NewQuery(filter.Contains(queryTagA))
	count := 0
	query.EachEntity(world, func(entry *donburi.Entry) {
		count++
		if entry.Archetype().Layout().HasComponent(queryTagA) == false {
			t.Errorf("PlayerTag should be in ent archetype")
		}
	})

	if count != 2 {
		t.Errorf("counter should be 2, but got %d", count)
	}
}

func TestQueryMultipleComponent(t *testing.T) {
	world := donburi.NewWorld()

	world.Create(queryTagA)
	world.Create(queryTagC)
	world.Create(queryTagA, queryTagB)

	query := donburi.NewQuery(filter.Contains(queryTagA, queryTagB))
	count := query.Count(world)
	if count != 1 {
		t.Errorf("counter should be 1, but got %d", count)
	}
}

func TestComplexQuery(t *testing.T) {
	createWorldFunc := func() donburi.World {
		world := donburi.NewWorld()

		world.Create(queryTagA)
		world.Create(queryTagC)
		world.Create(queryTagA, queryTagB)

		return world
	}

	var tests = []struct {
		filter        filter.LayoutFilter
		expectedCount int
	}{
		{filter.Not(filter.Contains(queryTagA)), 1},
		{filter.And(filter.Contains(queryTagA), filter.Not(filter.Contains(queryTagB))), 1},
		{filter.Or(filter.Contains(queryTagA), filter.Contains(queryTagC)), 3},
	}

	for _, tt := range tests {
		world := createWorldFunc()
		query := donburi.NewQuery(tt.filter)
		count := query.Count(world)
		if count != tt.expectedCount {
			t.Errorf("counter should be %d, but got %d", tt.expectedCount, count)
		}
	}
}

func TestFirstEntity(t *testing.T) {
	world := donburi.NewWorld()
	world.Create(queryTagA)
	world.Create(queryTagC)
	world.Create(queryTagA, queryTagB)

	// find first entity withqueryTagC
	query := donburi.NewQuery(filter.Contains(queryTagC))
	entry, ok := query.FirstEntity(world)
	if entry == nil || ok == false {
		t.Errorf("entry with queryTagC should not be nil")
	}

	entry.Remove()

	entry, ok = query.FirstEntity(world)
	if entry != nil || ok {
		t.Errorf("entry with queryTagC should be nil")
	}
}
