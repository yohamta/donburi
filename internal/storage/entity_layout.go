package storage

import "github.com/yohamta/donburi/internal/component"

// EntityLayout represents a layout of components.
type EntityLayout struct {
	components []*component.ComponentType
}

// NewEntityLayout creates a new entity layout.
func NewEntityLayout(components []*component.ComponentType) *EntityLayout {
	el := &EntityLayout{
		components: []*component.ComponentType{},
	}

	for _, ct := range components {
		el.RegisterComponent(ct)
	}

	return el
}

// Components returns the components of the layout.
func (el *EntityLayout) Components() []*component.ComponentType {
	return el.components
}

// RegisterComponent registers a component type to the layout.
func (el *EntityLayout) RegisterComponent(typeId *component.ComponentType) {
	el.components = append(el.components, typeId)
}

// HasComponent returns true if the layout has the given component type.
func (el *EntityLayout) HasComponent(componentType *component.ComponentType) bool {
	for _, ct := range el.components {
		if ct == componentType {
			return true
		}
	}
	return false
}
