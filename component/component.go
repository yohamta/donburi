package component

import (
	"reflect"
	"unsafe"
)

type (
	ComponentTypeId int
)

type IComponentType interface {
	Id() ComponentTypeId
	New() unsafe.Pointer
	Typ() reflect.Type
	Name() string
}
