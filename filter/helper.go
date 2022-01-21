package filter

import "github.com/yohamta/donburi/internal/component"

func containsComponent(components []*component.ComponentType, c *component.ComponentType) bool {
	for _, component := range components {
		if component == c {
			return true
		}
	}
	return false
}
