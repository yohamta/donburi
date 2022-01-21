package donburi_test

import (
	"testing"
	"unsafe"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/internal/entity"
)

type Vec2f struct {
	X float64
	Y float64
}

type TransformData struct {
	Position Vec2f
}

type VelocityData struct {
	Velocity Vec2f
}

var Transform = donburi.NewComponentType(TransformData{})
var Velocity = donburi.NewComponentType(VelocityData{})
var PlayerTag = donburi.NewTag()
var SecondaryTag = donburi.NewTag()
var EnemyTag = donburi.NewTag()

func TestEntry(t *testing.T) {
	nm := "TestEntry"
	world := donburi.NewWorld()

	player := world.Create(PlayerTag, Transform, Velocity)

	entry := world.Entry(player)
	if entry.Archetype().Layout().HasComponent(PlayerTag) == false {
		t.Errorf("%s: PlayerTag should be in player archetype", nm)
	}
}

func TestMutateComponent(t *testing.T) {
	nm := "TestMutateComponent"
	world := donburi.NewWorld()

	player := world.Create(PlayerTag, Transform, Velocity)
	enemy := world.Create(EnemyTag, Transform, Velocity)
	enemy2 := world.Create(EnemyTag, Transform, Velocity)

	p_entry := world.Entry(player)
	((*TransformData)(p_entry.Component(Transform))).Position.X = 10
	((*TransformData)(p_entry.Component(Transform))).Position.Y = 20

	e_entry := world.Entry(enemy)
	((*TransformData)(e_entry.Component(Transform))).Position.X = 30
	((*TransformData)(e_entry.Component(Transform))).Position.Y = 40

	e2_entry := world.Entry(enemy2)
	e2_entry.SetComponent(Transform, unsafe.Pointer(&TransformData{
		Position: Vec2f{40, 50},
	}))

	p_tr := (*TransformData)(p_entry.Component(Transform))
	if p_tr.Position.X != 10 || p_tr.Position.Y != 20 {
		t.Errorf("%s: Position should be (10, 20), but got (%f, %f)", nm, p_tr.Position.X, p_tr.Position.Y)
	}

	e_tr := (*TransformData)(e_entry.Component(Transform))
	if e_tr.Position.X != 30 || e_tr.Position.Y != 40 {
		t.Errorf("%s: Position should be (30, 40), but got (%f, %f)", nm, e_tr.Position.X, e_tr.Position.Y)
	}

	e2_tr := (*TransformData)(e2_entry.Component(Transform))
	if e2_tr.Position.X != 40 || e2_tr.Position.Y != 50 {
		t.Errorf("%s: Position should be (40, 50), but got (%f, %f)", nm, e2_tr.Position.X, e2_tr.Position.Y)
	}
}

func TestArchetype(t *testing.T) {
	nm := "TestArchetype"
	world := donburi.NewWorld()

	player := world.Create(PlayerTag, Transform, Velocity)

	entry := world.Entry(player)
	if entry.Archetype().Layout().HasComponent(PlayerTag) == false {
		t.Errorf("%s: PlayerTag should be in player archetype", nm)
	}
	if entry.Archetype().Layout().HasComponent(EnemyTag) == true {
		t.Errorf("%s: EnemyTag should not be in player archetype", nm)
	}
}

func TestAddComponent(t *testing.T) {
	nm := "TestAddComponent"
	world := donburi.NewWorld()

	entities := []entity.Entity{
		world.Create(PlayerTag, Transform),
		world.Create(PlayerTag, Transform),
		world.Create(PlayerTag, Transform),
	}

	entry := world.Entry(entities[1])
	old_arch := entry.Archetype()

	entry.AddComponent(Velocity, unsafe.Pointer(&VelocityData{
		Velocity: Vec2f{10, 20},
	}))
	entry.AddComponent(EnemyTag)

	new_arch := entry.Archetype()
	if len(new_arch.Layout().Components()) != 4 {
		t.Errorf("%s: Archetype should have 4 components", nm)
	}
	if len(old_arch.Entities()) != 2 {
		t.Errorf("%s: Old archetype should have 2 entities", nm)
	}
	if len(new_arch.Entities()) != 1 {
		t.Errorf("%s: Old archetype should have 1 entities", nm)
	}
}

func TestRemoveComponent(t *testing.T) {
	nm := "TestRemoveComponent"
	world := donburi.NewWorld()

	entities := []entity.Entity{
		world.Create(PlayerTag, Transform, Velocity),
		world.Create(PlayerTag, Transform, Velocity),
		world.Create(PlayerTag, Transform, Velocity),
	}

	entry := world.Entry(entities[1])
	old_arch := entry.Archetype()

	entry.RemoveComponent(Transform)

	new_arch := entry.Archetype()
	if len(new_arch.Layout().Components()) != 2 {
		t.Errorf("%s: Archetype should have 2 components", nm)
	}
	if len(old_arch.Entities()) != 2 {
		t.Errorf("%s: Old archetype should have 2 entities", nm)
	}
	if len(new_arch.Entities()) != 1 {
		t.Errorf("%s: Old archetype should have 1 entities", nm)
	}
}

func TestDeleteEntity(t *testing.T) {
	var tests = []struct {
		name string
		del  []int
		want int
	}{
		{"Delete first", []int{0}, 2},
		{"Delete all", []int{0, 1, 2}, 0},
		{"Delete last", []int{2}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := donburi.NewWorld()
			entities := []donburi.Entity{
				world.Create(PlayerTag, Transform, Velocity),
				world.Create(PlayerTag, Transform, Velocity),
				world.Create(PlayerTag, Transform, Velocity),
			}
			for _, del := range tt.del {
				world.Remove(entities[del])
				if world.Valid(entities[del]) {
					t.Errorf("%s: Entity should be deleted", tt.name)
				}
			}
			if world.Len() != tt.want {
				t.Errorf("%s: want %d, got %d", tt.name, tt.want, world.Len())
			}
		})
	}

}
