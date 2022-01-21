package filter

import (
	"github.com/yohamta/donburi/internal/component"
)

type exact struct {
	components []*component.ComponentType
}

// Exact matches layouts that contain exactly the same components specified.
func Exact(components []*component.ComponentType) LayoutFilter {
	return exact{
		components: components,
	}
}

func (f exact) MatchesLayout(components []*component.ComponentType) bool {
	if len(components) != len(f.components) {
		return false
	}
	for _, componentType := range components {
		if !containsComponent(f.components, componentType) {
			return false
		}
	}
	return true
}
