package query_test

import (
	"testing"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/query"
)

var PlayerTag = donburi.NewTag()
var SecondaryTag = donburi.NewTag()
var EnemyTag = donburi.NewTag()

func TestQuery(t *testing.T) {
	nm := "TestQuery"

	world := donburi.NewWorld()
	world.Create(PlayerTag)
	world.Create(EnemyTag)
	world.Create(PlayerTag, SecondaryTag)

	query := query.NewQuery(filter.Contains(PlayerTag))
	count := 0
	query.EachEntity(world, func(entry *donburi.Entry) {
		count++
		if entry.Archetype().Layout().HasComponent(PlayerTag) == false {
			t.Errorf("%s: PlayerTag should be in player archetype", nm)
		}
	})

	if count != 2 {
		t.Errorf("%s: counter should be 2, but got %d", nm, count)
	}
}

func TestQueryMultipleComponent(t *testing.T) {
	nm := "TestQueryEntities"

	world := donburi.NewWorld()
	world.Create(PlayerTag)
	world.Create(EnemyTag)
	world.Create(PlayerTag, SecondaryTag)

	query := query.NewQuery(filter.Contains(PlayerTag, SecondaryTag))

	count := 0
	query.EachEntity(world, func(_ *donburi.Entry) {
		count++
	})

	if count != 1 {
		t.Errorf("%s: counter should be 1, but got %d", nm, count)
	}
}

func TestComplexQuery(t *testing.T) {
	var tests = []struct {
		name   string
		filter filter.LayoutFilter
		want   int
	}{
		{"Not filter", filter.Not(filter.Contains(PlayerTag)), 1},
		{"And filter", filter.And(filter.Contains(PlayerTag), filter.Not(filter.Contains(SecondaryTag))), 1},
		{"Or filter", filter.Or(filter.Contains(PlayerTag), filter.Contains(EnemyTag)), 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := donburi.NewWorld()
			world.Create(PlayerTag)
			world.Create(EnemyTag)
			world.Create(PlayerTag, SecondaryTag)

			query := query.NewQuery(tt.filter)
			count := 0
			query.EachEntity(world, func(entry *donburi.Entry) {
				count++
			})

			if count != tt.want {
				t.Errorf("%s: counter should be %d, but got %d", tt.name, tt.want, count)
			}

		})
	}

}
