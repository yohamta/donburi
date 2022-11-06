package filter

import "github.com/yohamta/donburi/internal/component"

type contains struct {
	components []component.IComponentType
}

// Contains matches layouts that contains all the components specified.
func Contains(components ...component.IComponentType) LayoutFilter {
	return &contains{components: components}
}

func (f *contains) MatchesLayout(components []component.IComponentType) bool {
	for _, componentType := range f.components {
		if !containsComponent(components, componentType) {
			return false
		}
	}
	return true
}
