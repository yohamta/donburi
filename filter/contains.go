package filter

import "github.com/yohamta/donburi/internal/component"

type contains struct {
	components []*component.ComponentType
}

// Contains matches layouts that contains all the components specified.
func Contains(components ...*component.ComponentType) LayoutFilter {
	return &contains{components: components}
}

func (f *contains) MatchesLayout(components []*component.ComponentType) bool {
	for _, componentType := range f.components {
		if !containsComponent(components, componentType) {
			return false
		}
	}
	return true
}
