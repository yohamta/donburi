package storage

import (
	"testing"

	"github.com/yohamta/donburi/component"
)

func TestComponents(t *testing.T) {
	type ComponentData struct {
		ID string
	}
	var (
		ca = NewMockComponentType(ComponentData{}, nil)
		cb = NewMockComponentType(ComponentData{}, nil)
	)

	components := NewComponents()

	tests := []*struct {
		layout  *Layout
		archIdx ArchetypeIndex
		compIdx ComponentIndex
		ID      string
	}{
		{
			NewLayout([]component.IComponentType{ca}),
			0,
			0,
			"a",
		},
		{
			NewLayout([]component.IComponentType{ca, cb}),
			1,
			1,
			"b",
		},
	}

	for _, tt := range tests {
		tt.compIdx = components.PushComponents(tt.layout.Components(), tt.archIdx)
	}

	for _, tt := range tests {
		for _, comp := range tt.layout.Components() {
			st := components.Storage(comp)
			if !st.Contains(tt.archIdx, tt.compIdx) {
				t.Errorf("storage should contain the component at %d, %d", tt.archIdx, tt.compIdx)
			}
			dat := (*ComponentData)(st.Component(tt.archIdx, tt.compIdx))
			dat.ID = tt.ID
		}
	}

	target := tests[0]
	storage := components.Storage(ca)

	var srcArchIdx ArchetypeIndex = target.archIdx
	var dstArchIdx ArchetypeIndex = 1

	storage.MoveComponent(srcArchIdx, target.compIdx, dstArchIdx)
	components.Move(srcArchIdx, dstArchIdx)

	if storage.Contains(srcArchIdx, target.compIdx) {
		t.Errorf("storage should not contain the component at %d, %d", target.archIdx, target.compIdx)
	}
	if components.componentIndices[srcArchIdx] != -1 {
		t.Errorf("component index should be -1 at %d but %d", srcArchIdx, components.componentIndices[srcArchIdx])
	}

	newCompIdx := components.componentIndices[dstArchIdx]
	if !storage.Contains(dstArchIdx, newCompIdx) {
		t.Errorf("storage should contain the component at %d, %d", dstArchIdx, target.compIdx)
	}

	dat := (*ComponentData)(storage.Component(dstArchIdx, newCompIdx))
	if dat.ID != target.ID {
		t.Errorf("component should have ID %s", target.ID)
	}
}
