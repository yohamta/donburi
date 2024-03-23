package donburi_test

import (
	"testing"

	"github.com/yohamta/donburi"
)

func TestSetValue(t *testing.T) {
	var (
		transform = donburi.NewComponentType[transformData]()
	)
	world := donburi.NewWorld()
	a := world.Create(transform)
	entryA := world.Entry(a)

	trData := transformData{
		Position: vec2f{10, 20},
	}
	donburi.SetValue(entryA, transform, trData)
	got := donburi.GetValue[transformData](entryA, transform)

	if got != trData {
		t.Errorf("got: %v, want: %v", got, trData)
	}
}

func TestGetComponents(t *testing.T) {
	var (
		transform = donburi.NewComponentType[transformData]()
		velocity  = donburi.NewComponentType[velocityData]()
		tag       = donburi.NewTag().SetName("tag")
	)

	world := donburi.NewWorld()
	a := world.Create(transform, velocity, tag)
	entryA := world.Entry(a)

	trData := transformData{
		Position: vec2f{10, 20},
	}
	veData := velocityData{
		Velocity: vec2f{30, 40},
	}

	donburi.SetValue(entryA, transform, trData)
	donburi.SetValue(entryA, velocity, veData)

	gots := donburi.GetComponents(entryA)
	wants := []interface{}{trData, veData, struct{}{}}

	if len(gots) != len(wants) {
		t.Fatalf("got: %v, want: %v", gots, wants)
	}

	for i, got := range gots {
		if got != wants[i] {
			t.Errorf("got: %v, want: %v", got, wants[i])
		}
	}
}
