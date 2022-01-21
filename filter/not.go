package filter

import "github.com/yohamta/donburi/internal/component"

type not struct {
	filter LayoutFilter
}

func Not(filter LayoutFilter) LayoutFilter {
	return &not{filter: filter}
}

func (f *not) MatchesLayout(components []*component.ComponentType) bool {
	return !f.filter.MatchesLayout(components)
}
