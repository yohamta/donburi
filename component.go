package donburi

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/filter"
)

// IComponentType is an interface for component types.
type IComponentType = component.IComponentType

// NewComponentType creates a new component type.
// The function is used to create a new component of the type.
// It receives a function that returns a pointer to a new component.
// The first argument is a default value of the component.
func NewComponentType[T any](opts ...interface{}) *ComponentType[T] {
	var t T
	if len(opts) == 0 {
		return newComponentType(t, nil)
	}
	return newComponentType(t, opts[0])
}

// CompnentType represents a type of component. It is used to identify
// a component when getting or setting components of an entity.
type ComponentType[T any] struct {
	id         component.ComponentTypeId
	typ        reflect.Type
	name       string
	defaultVal interface{}
	query      *Query
}

// Typ returns the reflect.Type of the ComponentType.
func (c *ComponentType[T]) Typ() reflect.Type {
	return c.typ
}

// Get returns component data from the entry.
func (c *ComponentType[T]) Get(entry *Entry) *T {
	return (*T)(entry.Component(c))
}

// GetValue returns value of the component from the entry.
func (c *ComponentType[T]) GetValue(entry *Entry) T {
	return *Get[T](entry, c)
}

// Set sets component data to the entry.
func (c *ComponentType[T]) Set(entry *Entry, component *T) {
	entry.SetComponent(c, unsafe.Pointer(component))
}

// Each iterates over the entities that have the component.
func (c *ComponentType[T]) Each(w World, callback func(*Entry)) {
	c.query.Each(w, callback)
}

// deprecated: use Each instead
func (c *ComponentType[T]) EachEntity(w World, callback func(*Entry)) {
	c.Each(w, callback)
}

// First returns the first entity that has the component.
func (c *ComponentType[T]) First(w World) (*Entry, bool) {
	return c.query.First(w)
}

// deprecated: use First instead
func (c *ComponentType[T]) FirstEntity(w World) (*Entry, bool) {
	return c.First(w)
}

// MustFirst returns the first entity that has the component or panics.
func (c *ComponentType[T]) MustFirst(w World) *Entry {
	e, ok := c.query.First(w)
	if !ok {
		panic(fmt.Sprintf("no entity has the component %s", c.name))
	}

	return e
}

// deprecated: use MustFirst instead
func (c *ComponentType[T]) MustFirstEntity(w World) *Entry {
	return c.MustFirst(w)
}

// SetValue sets the value of the component.
func (c *ComponentType[T]) SetValue(entry *Entry, value T) {
	comp := c.Get(entry)
	*comp = value
}

// String returns the component type name.
func (c *ComponentType[T]) String() string {
	return c.name
}

// SetName sets the component type name.
func (c *ComponentType[T]) SetName(name string) *ComponentType[T] {
	c.name = name
	return c
}

// Name returns the component type name.
func (c *ComponentType[T]) Name() string {
	return c.name
}

// Id returns the component type id.
func (c *ComponentType[T]) Id() component.ComponentTypeId {
	return c.id
}

func (c *ComponentType[T]) New() unsafe.Pointer {
	val := reflect.New(c.typ)
	v := reflect.Indirect(val)
	ptr := unsafe.Pointer(v.UnsafeAddr())
	if c.defaultVal != nil {
		c.setDefaultVal(ptr)
	}
	return ptr
}

func (c *ComponentType[T]) setDefaultVal(ptr unsafe.Pointer) {
	v := reflect.Indirect(reflect.ValueOf(c.defaultVal))
	reflect.NewAt(c.typ, ptr).Elem().Set(v)
}

func (c *ComponentType[T]) validateDefaultVal() {
	if !reflect.TypeOf(c.defaultVal).AssignableTo(c.typ) {
		err := fmt.Sprintf("default value is not assignable to component type: %s", c.name)
		panic(err)
	}
}

var nextComponentTypeId component.ComponentTypeId = 1

// NewComponentType creates a new component type.
// The argument is a struct that represents a data of the component.
func newComponentType[T any](s T, defaultVal interface{}) *ComponentType[T] {
	componentType := &ComponentType[T]{
		id:         nextComponentTypeId,
		typ:        reflect.TypeOf(s),
		name:       reflect.TypeOf(s).Name(),
		defaultVal: defaultVal,
	}
	componentType.query = NewQuery(filter.Contains(componentType))
	if defaultVal != nil {
		componentType.validateDefaultVal()
	}
	nextComponentTypeId++
	return componentType
}
