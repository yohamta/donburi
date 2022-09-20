package storage

import (
	"testing"

	"github.com/yohamta/donburi/internal/component"
)

type testComponentData struct {
}

var testComponentType = component.NewComponentType(testComponentData{})

func TestLayout(t *testing.T) {
	components := []*component.ComponentType{testComponentType}
	layout := NewLayout(components)

	if layout.HasComponent(testComponentType) == false {
		t.Errorf("layout should have the component type")
	}
}
