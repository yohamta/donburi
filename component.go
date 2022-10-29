package donburi

import "github.com/yohamta/donburi/internal/component"

// ComponentType represents a component type.
// It is used to add components to entities, and to filter entities based on their components.
// It contains a function that returns a pointer to a new component.
type ComponentType = component.ComponentType

// NewComponentType creates a new component type.
// The function is used to create a new component of the type.
// It receives a function that returns a pointer to a new component.
func NewComponentType[T any](opts ...interface{}) *ComponentType {
	var t T
	if len(opts) == 0 {
		return component.NewComponentType(t, nil)
	}
	return component.NewComponentType(t, opts[0])
}
