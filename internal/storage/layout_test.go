package storage

import (
	"testing"

	"github.com/yohamta/donburi/internal/component"
)

func TestLayout(t *testing.T) {
	compType := component.NewComponentType(struct{}{}, nil)
	components := []*component.ComponentType{compType}
	layout := NewLayout(components)

	if layout.HasComponent(compType) == false {
		t.Errorf("layout should have the component type")
	}
}
