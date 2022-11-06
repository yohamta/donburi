package filter

import (
	"github.com/yohamta/donburi/internal/component"
)

type exact struct {
	components []component.IComponentType
}

// Exact matches layouts that contain exactly the same components specified.
func Exact(components []component.IComponentType) LayoutFilter {
	return exact{
		components: components,
	}
}

func (f exact) MatchesLayout(components []component.IComponentType) bool {
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
