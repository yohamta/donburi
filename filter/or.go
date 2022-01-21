package filter

import "github.com/yohamta/donburi/internal/component"

type or struct {
	filters []LayoutFilter
}

func Or(filters ...LayoutFilter) LayoutFilter {
	return &or{filters: filters}
}

func (f *or) MatchesLayout(components []*component.ComponentType) bool {
	for _, filter := range f.filters {
		if filter.MatchesLayout(components) {
			return true
		}
	}
	return false
}
