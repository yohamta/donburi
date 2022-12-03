package storage

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/yohamta/donburi/component"
)

var (
	nextMockComponentTypeId component.ComponentTypeId = 1
)

type MockComponentType[T any] struct {
	id         component.ComponentTypeId
	typ        reflect.Type
	defaultVal interface{}
}

func NewMockComponentType[T any](t T, defaultVal interface{}) *MockComponentType[T] {
	m := &MockComponentType[T]{
		id:         nextMockComponentTypeId,
		typ:        reflect.TypeOf(t),
		defaultVal: defaultVal,
	}
	nextMockComponentTypeId++
	return m
}

func (m *MockComponentType[T]) Id() component.ComponentTypeId {
	return m.id
}

func (m *MockComponentType[T]) New() unsafe.Pointer {
	val := reflect.New(m.typ)
	v := reflect.Indirect(val)
	ptr := unsafe.Pointer(v.UnsafeAddr())
	if m.defaultVal != nil {
		m.setDefaultVal(ptr)
	}
	return ptr
}

func (m *MockComponentType[T]) setDefaultVal(ptr unsafe.Pointer) {
	v := reflect.Indirect(reflect.ValueOf(m.defaultVal))
	reflect.NewAt(m.typ, ptr).Elem().Set(v)
}

func (m *MockComponentType[T]) Name() string {
	return fmt.Sprintf("%s[%s]", reflect.TypeOf(m).Name(), m.typ.Name())
}
