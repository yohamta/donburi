package storage

import (
	"testing"
	"unsafe"
)

func TestStorage(t *testing.T) {
	type Component struct{ ID string }
	var (
		componentType = NewMockComponentType[any](Component{}, nil)
	)

	st := NewStorage()

	tests := []struct {
		ID       string
		expected string
	}{
		{ID: "a", expected: "a"},
		{ID: "b", expected: "b"},
		{ID: "c", expected: "c"},
	}

	var archIdx ArchetypeIndex = 0
	var compIdx ComponentIndex = 0
	for _, tt := range tests {
		st.PushComponent(componentType, archIdx)
		st.SetComponent(archIdx, compIdx, unsafe.Pointer(&Component{ID: tt.ID}))
		compIdx++
	}

	compIdx = 0
	for _, tt := range tests {
		if c := (*Component)(st.Component(archIdx, compIdx)); c.ID != tt.expected {
			t.Errorf("component should have ID %s", tt.expected)
		}
		compIdx++
	}

	// SwapRemove component
	removed := (*Component)(st.SwapRemove(archIdx, 1))
	if removed == nil {
		t.Fatalf("removed component should not be nil")
	}
	if removed.ID != "b" {
		t.Errorf("removed component should have ID b")
	}

	tests2 := []struct {
		archIdx    ArchetypeIndex
		cmpIdx     ComponentIndex
		epxectedID string
	}{
		{archIdx: 0, cmpIdx: 0, epxectedID: "a"},
		{archIdx: 0, cmpIdx: 1, epxectedID: "c"},
	}

	for _, tt := range tests2 {
		if c := (*Component)(st.Component(tt.archIdx, tt.cmpIdx)); c.ID != tt.epxectedID {
			t.Errorf("component should have ID %s but got %s", tt.epxectedID, c.ID)
		}
		compIdx++
	}
}
