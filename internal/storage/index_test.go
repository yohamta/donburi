package storage

import (
	"testing"

	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/internal/component"
)

func TestIndex(t *testing.T) {
	var (
		ca = component.NewComponentType(struct{}{})
		cb = component.NewComponentType(struct{}{})
		cc = component.NewComponentType(struct{}{})
	)

	index := NewIndex()

	layoutA := NewLayout([]*component.ComponentType{ca})
	layoutB := NewLayout([]*component.ComponentType{ca, cb})

	index.Push(layoutA)
	index.Push(layoutB)

	tests := []struct {
		filter   filter.LayoutFilter
		expected int
	}{
		{

			filter:   filter.Contains(ca),
			expected: 2,
		},
		{

			filter:   filter.Contains(cb),
			expected: 1,
		},
		{

			filter:   filter.Contains(cc),
			expected: 0,
		},
	}

	for _, tt := range tests {
		it := index.Search(tt.filter)
		if len(it.values) != tt.expected {
			t.Errorf("index should have %d archetypes", tt.expected)
		}
		if it.current != 0 && it.HasNext() {
			t.Errorf("index should have 0 as current")
		}
		if tt.expected == 0 && it.HasNext() {
			t.Errorf("index should not have next")
		}
	}
}
