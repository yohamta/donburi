package donburi_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/internal/storage"
)

type vec2f struct {
	X float64
	Y float64
}

type transformData struct {
	Position vec2f
}

type velocityData struct {
	Velocity vec2f
}

var (
	transform = donburi.NewComponentType[transformData]()
	velocity  = donburi.NewComponentType[velocityData]()
	tagA      = donburi.NewTag()
	tagB      = donburi.NewTag()
)

func TestEntry(t *testing.T) {
	world := donburi.NewWorld()
	ent := world.Create(tagA, transform, velocity)
	entry := world.Entry(ent)

	if !entry.HasComponent(tagA) {
		t.Fatalf("TagA should be in ent archetype")
	}
}

func TestMutateComponent(t *testing.T) {
	world := donburi.NewWorld()

	a := world.Create(tagA, transform, velocity)
	b := world.Create(tagB, transform, velocity)
	c := world.Create(tagB, transform, velocity)

	entryA := world.Entry(a)
	donburi.Get[transformData](entryA, transform).Position.X = 10
	donburi.Get[transformData](entryA, transform).Position.Y = 20

	entryB := world.Entry(b)
	donburi.Get[transformData](entryB, transform).Position.X = 30
	donburi.Get[transformData](entryB, transform).Position.Y = 40

	entryC := world.Entry(c)
	tr := &transformData{Position: vec2f{40, 50}}
	donburi.Add(entryC, transform, tr)

	tests := []struct {
		entry    *donburi.Entry
		expected *vec2f
	}{
		{entryA, &vec2f{10, 20}},
		{entryB, &vec2f{30, 40}},
		{entryC, &vec2f{40, 50}},
	}

	for _, tt := range tests {
		tf := donburi.Get[transformData](tt.entry, transform)
		if tf.Position.X != tt.expected.X {
			t.Errorf("X should be %f, but %f", tt.expected.X, tf.Position.X)
		}
		if tf.Position.Y != tt.expected.Y {
			t.Errorf("Y should be %f, but %f", tt.expected.Y, tf.Position.Y)
		}
	}
}

func TestArchetype(t *testing.T) {
	world := donburi.NewWorld()
	entity := world.Create(tagA, transform, velocity)
	entry := world.Entry(entity)

	if entry.HasComponent(tagA) == false {
		t.Errorf("TagA should be in the archetype")
	}
	if entry.HasComponent(tagB) == true {
		t.Errorf("TagB should not be in the archetype")
	}
}

func TestAddComponent(t *testing.T) {
	world := donburi.NewWorld()

	entities := []storage.Entity{
		world.Create(tagA, transform),
		world.Create(tagA, transform),
		world.Create(tagA, transform),
	}

	entry := world.Entry(entities[1])
	archtype := entry.Archetype()

	vd := &velocityData{Velocity: vec2f{10, 20}}
	donburi.Add(entry, velocity, vd)
	entry.AddComponent(tagB)

	newArchtype := entry.Archetype()
	if len(newArchtype.Layout().Components()) != 4 {
		t.Errorf("New Archetype should have 4 components")
	}
	if len(archtype.Entities()) != 2 {
		t.Errorf("Old archetype should have 2 entities")
	}
	if len(newArchtype.Entities()) != 1 {
		t.Errorf("New archetype should have 1 entities")
	}
}

func TestRemoveComponent(t *testing.T) {
	world := donburi.NewWorld()

	entities := []storage.Entity{
		world.Create(tagA, transform, velocity),
		world.Create(tagA, transform, velocity),
		world.Create(tagA, transform, velocity),
	}

	entry := world.Entry(entities[1])
	archtype := entry.Archetype()

	entry.RemoveComponent(transform)

	newArchtype := entry.Archetype()
	if len(newArchtype.Layout().Components()) != 2 {
		t.Errorf("Archetype should have 2 components")
	}
	if len(archtype.Entities()) != 2 {
		t.Errorf("Old archetype should have 2 entities")
	}
	if len(newArchtype.Entities()) != 1 {
		t.Errorf("New archetype should have 1 entities")
	}
}

func TestDeleteEntity(t *testing.T) {
	createWorldFunc := func() (donburi.World, []*donburi.Entry) {
		world := donburi.NewWorld()
		entities := []*donburi.Entry{
			world.Entry(world.Create(tagA, transform, velocity)),
			world.Entry(world.Create(tagA, transform, velocity)),
			world.Entry(world.Create(tagA, transform, velocity)),
		}
		return world, entities
	}

	var tests = []struct {
		DeleteIndex         []int
		expectedEntityCount int
	}{
		{[]int{0}, 2},
		{[]int{0, 1, 2}, 0},
		{[]int{2}, 2},
	}

	for _, tt := range tests {
		world, entities := createWorldFunc()

		for _, ent := range entities {
			if !world.Valid(ent.Entity()) {
				t.Errorf("Entity should be valid")
			}
			if !ent.Valid() {
				t.Errorf("Entry should be valid")
			}
		}

		for _, del := range tt.DeleteIndex {
			world.Remove(entities[del].Entity())
			if world.Valid(entities[del].Entity()) {
				t.Errorf("Entity should be invalid")
			}
			if entities[del].Valid() {
				t.Errorf("Entry should be invalid")
			}
		}

		if world.Len() != tt.expectedEntityCount {
			t.Errorf("World should have %d entities", tt.expectedEntityCount)
		}
	}
}

func TestRemoveAndCreateEntity(t *testing.T) {
	world := donburi.NewWorld()

	entityA := world.Create(tagA)

	world.Remove(entityA)
	require.False(t, world.Valid(entityA))

	entityB := world.Create(tagA)

	query := donburi.NewQuery(filter.Contains(tagA))
	entry, ok := query.First(world)
	if !ok {
		t.Fatalf("Entity should be found")
	}
	if entry.Entity() != entityB {
		t.Errorf("Entity should be %d, but %d", entityB, entry.Entity())
	}
	if !entry.HasComponent(tagA) {
		t.Errorf("TagA should be in the archetype")
	}
}

type archeTest struct {
	name          string
	archetype     *storage.Archetype
	expectedCount int
}

func TestCreateEntityAndExtend(t *testing.T) {
	testFunc := func(tests []archeTest) {
		t.Helper()
		for _, tt := range tests {
			if len(tt.archetype.Entities()) != tt.expectedCount {
				t.Errorf("%s archetype should have %d entities", tt.name, tt.expectedCount)
			}
		}
	}
	world := donburi.NewWorld()

	entity := world.Create(velocity)
	entry := world.Entry(entity)

	oldArchtype := entry.Archetype()

	testFunc([]archeTest{
		{"old", oldArchtype, 1},
	})

	// Add new component
	donburi.Add(entry, transform, &transformData{})
	newArchtype := entry.Archetype()

	testFunc([]archeTest{
		{"old", oldArchtype, 0},
		{"new", newArchtype, 1},
	})

	// Create another entity
	anotherEntity := world.Create(velocity)
	anotherEntry := world.Entry(anotherEntity)

	testFunc([]archeTest{
		{"old", oldArchtype, 1},
		{"new", newArchtype, 1},
	})

	donburi.Add(anotherEntry, transform, &transformData{})

	testFunc([]archeTest{
		{"old", oldArchtype, 0},
		{"new", newArchtype, 2},
	})
}

func TestComponentDefaultVal(t *testing.T) {
	type ComponentData struct {
		Val int
	}
	defVal := ComponentData{Val: 10}
	component := donburi.NewComponentType[ComponentData](defVal)

	world := donburi.NewWorld()

	entry := world.Entry(world.Create(component))
	val := donburi.Get[ComponentData](entry, component)

	if val.Val != defVal.Val {
		t.Errorf("Default value should be %d, but %d", defVal.Val, val.Val)
	}
}
