package storage

import (
	"bytes"
	"fmt"

	"github.com/yohamta/donburi/internal/component"
)

// Layout represents a layout of components.
type Layout struct {
	componentTypes []*component.ComponentType
}

// NewLayout creates a new entity layout.
func NewLayout(components []*component.ComponentType) *Layout {
	layout := &Layout{
		componentTypes: []*component.ComponentType{},
	}

	for _, ct := range components {
		layout.AddComponent(ct)
	}

	return layout
}

// Components returns the components of the layout.
func (l *Layout) Components() []*component.ComponentType {
	return l.componentTypes
}

// AddComponent registers a component type to the layout.
func (l *Layout) AddComponent(componentType *component.ComponentType) {
	l.componentTypes = append(l.componentTypes, componentType)
}

// HasComponent returns true if the layout has the given component type.
func (l *Layout) HasComponent(componentType *component.ComponentType) bool {
	for _, ct := range l.componentTypes {
		if ct == componentType {
			return true
		}
	}
	return false
}

func (l *Layout) String() string {
	var out bytes.Buffer
	out.WriteString("Layout: {")
	for i, ct := range l.componentTypes {
		if i != 0 {
			out.WriteString(", ")
		}
		out.WriteString(fmt.Sprintf("%s", ct))
	}
	out.WriteString("}")
	return out.String()
}
