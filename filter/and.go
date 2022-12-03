package filter

import "github.com/yohamta/donburi/component"

type and struct {
	filters []LayoutFilter
}

func And(filters ...LayoutFilter) LayoutFilter {
	return &and{filters: filters}
}

func (f *and) MatchesLayout(components []component.IComponentType) bool {
	for _, filter := range f.filters {
		if !filter.MatchesLayout(components) {
			return false
		}
	}
	return true
}
