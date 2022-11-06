package donburi

import "github.com/yohamta/donburi/internal/component"

type ComponentType = component.IComponentType

// NewComponentType creates a new component type.
// The function is used to create a new component of the type.
// It receives a function that returns a pointer to a new component.
// The first argument is a default value of the component.
func NewComponentType[T any](opts ...interface{}) *component.ComponentType[T] {
	var t T
	if len(opts) == 0 {
		return component.NewComponentType(t, nil)
	}
	return component.NewComponentType(t, opts[0])
}
