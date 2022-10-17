package component

import (
	"reflect"
	"unsafe"
)

type ComponentTypeId int

// CompnentType represents a type of component. It is used to identify
// a component when getting or setting components of an entity.
type ComponentType struct {
	id   ComponentTypeId
	typ  reflect.Type
	name string
}

var nextComponentTypeId ComponentTypeId = 1

// NewComponentType creates a new component type.
// The argument is a struct that represents a data of the component.
func NewComponentType(s interface{}) *ComponentType {
	componentType := &ComponentType{
		id:   nextComponentTypeId,
		typ:  reflect.TypeOf(s),
		name: reflect.TypeOf(s).Name(),
	}
	nextComponentTypeId++
	return componentType
}

// String returns the component type name.
func (c *ComponentType) String() string {
	return c.name
}

// SetName sets the component type name.
func (c *ComponentType) SetName(name string) *ComponentType {
	c.name = name
	return c
}

// Name returns the component type name.
func (c *ComponentType) Name() string {
	return c.name
}

// Id returns the component type id.
func (c *ComponentType) Id() ComponentTypeId {
	return c.id
}

func (c *ComponentType) New() unsafe.Pointer {
	val := reflect.New(c.typ)
	v := reflect.Indirect(val)
	return unsafe.Pointer(v.UnsafeAddr())
}
