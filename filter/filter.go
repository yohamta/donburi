package filter

import "github.com/yohamta/donburi/internal/component"

// LayoutFilter is a filter that filters entities based on their components.
type LayoutFilter interface {

	// MatchesLayout returns true if the entity matches the filter.
	MatchesLayout(components []component.IComponentType) bool
}
