package component

import (
	"fmt"

	reflect "github.com/goccy/go-reflect"
)

type ComponentTypeId int

// CompnentType represents a type of component. It is used to identify
// a component when getting or setting components of an entity.
type ComponentType struct {
	id   ComponentTypeId
	zero reflect.Value
}

var nextComponentTypeId ComponentTypeId = 1

// NewComponentType creates a new component type.
// The argument is a struct that represents a data of the component.
func NewComponentType(s interface{}) *ComponentType {
	if err := validate(s); err != nil {
		panic(err)
	}
	componentType := &ComponentType{
		id:   nextComponentTypeId,
		zero: reflect.Zero(reflect.TypeOf(s)),
	}
	nextComponentTypeId++
	return componentType
}

// Id returns the component type id.
func (c *ComponentType) Id() ComponentTypeId {
	return c.id
}

// New creates a new component value of the stored type.
func (c *ComponentType) New() reflect.Value {
	return c.zero
}

func validate(s interface{}) error {
	typ := reflect.TypeOf(s)
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("component must be struct, but %v", typ.Kind())
	}
	return nil
}
