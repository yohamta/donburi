package storage

import (
	"testing"

	"github.com/yohamta/donburi/internal/component"
)

type (
	componentA = struct{}
	componentB = struct{}
)

func TestMatchesLayout(t *testing.T) {
	var (
		ca = component.NewMockComponentType(componentA{}, nil)
		cb = component.NewMockComponentType(componentB{}, nil)
	)

	cmps := []component.IComponentType{ca, cb}
	archetype := NewArchetype(0, NewLayout(cmps))
	if !archetype.LayoutMatches(cmps) {
		t.Errorf("archetype should match the layout")
	}
}

func TestPushEntity(t *testing.T) {
	var (
		ca = component.NewMockComponentType(struct{}{}, nil)
		cb = component.NewMockComponentType(struct{}{}, nil)
	)

	cmps := []component.IComponentType{ca, cb}
	archetype := NewArchetype(0, NewLayout(cmps))

	archetype.PushEntity(0)
	archetype.PushEntity(1)
	archetype.PushEntity(2)

	if len(archetype.Entities()) != 3 {
		t.Errorf("archetype should have 3 entities")
	}

	archetype.SwapRemove(1)
	if len(archetype.Entities()) != 2 {
		t.Errorf("archetype should have 2 entities")
	}

	expected := []int{0, 2}
	for i, entity := range archetype.Entities() {
		if int(entity) != expected[i] {
			t.Errorf("archetype should have entity %d", expected[i])
		}
	}
}
