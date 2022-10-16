package query_test

import (
	"testing"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

var (
	tagA = donburi.NewTag()
	tagB = donburi.NewTag()
	tagC = donburi.NewTag()
)

func TestQuery(t *testing.T) {
	world := donburi.NewWorld()
	world.Create(tagA)
	world.Create(tagC)
	world.Create(tagA, tagB)

	query := query.NewQuery(filter.Contains(tagA))
	count := 0
	query.EachEntity(world, func(entry *donburi.Entry) {
		count++
		if entry.Archetype().Layout().HasComponent(tagA) == false {
			t.Errorf("PlayerTag should be in ent archetype")
		}
	})

	if count != 2 {
		t.Errorf("counter should be 2, but got %d", count)
	}
}

func TestQueryMultipleComponent(t *testing.T) {
	world := donburi.NewWorld()

	world.Create(tagA)
	world.Create(tagC)
	world.Create(tagA, tagB)

	query := query.NewQuery(filter.Contains(tagA, tagB))
	count := query.Count(world)
	if count != 1 {
		t.Errorf("counter should be 1, but got %d", count)
	}
}

func TestComplexQuery(t *testing.T) {
	createWorldFunc := func() donburi.World {
		world := donburi.NewWorld()

		world.Create(tagA)
		world.Create(tagC)
		world.Create(tagA, tagB)

		return world
	}

	var tests = []struct {
		filter        filter.LayoutFilter
		expectedCount int
	}{
		{filter.Not(filter.Contains(tagA)), 1},
		{filter.And(filter.Contains(tagA), filter.Not(filter.Contains(tagB))), 1},
		{filter.Or(filter.Contains(tagA), filter.Contains(tagC)), 3},
	}

	for _, tt := range tests {
		world := createWorldFunc()
		query := query.NewQuery(tt.filter)
		count := query.Count(world)
		if count != tt.expectedCount {
			t.Errorf("counter should be %d, but got %d", tt.expectedCount, count)
		}
	}
}

func TestFirstEntity(t *testing.T) {
	world := donburi.NewWorld()
	world.Create(tagA)
	world.Create(tagC)
	world.Create(tagA, tagB)

	// find first entity with tagC
	query := query.NewQuery(filter.Contains(tagC))
	entry, ok := query.FirstEntity(world)
	if entry == nil || ok == false {
		t.Errorf("entry with tagC should not be nil")
	}

	entry.Remove()

	entry, ok = query.FirstEntity(world)
	if entry != nil || ok {
		t.Errorf("entry with tagC should be nil")
	}
}
