package component

import (
	"unsafe"
)

type (
	ComponentTypeId int
)

type IComponentType interface {
	Id() ComponentTypeId
	New() unsafe.Pointer
	Name() string
}
