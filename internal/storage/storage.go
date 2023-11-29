package storage

import (
	"unsafe"

	"github.com/yohamta/donburi/component"
)

// Storage is a structure that stores the pointer to data of each component.
// It stores the pointers in the two dimensional slice.
// First dimension is the archetype index.
// Second dimension is the component index.
// The component index is used to access the component data in the archetype.
type Storage struct {
	storages [][]unsafe.Pointer
}

// NewStorage creates a new empty structure that stores the pointer to data of each component.
func NewStorage() *Storage {
	return &Storage{
		storages: make([][]unsafe.Pointer, 256),
	}
}

// PushComponent stores the new data of the component in the archetype.
func (cs *Storage) PushComponent(component component.IComponentType, archetypeIndex ArchetypeIndex) {
	if len(cs.storages) <= int(archetypeIndex) {
		cs.ensureCapacity(archetypeIndex)
	}
	if v := cs.storages[archetypeIndex]; v == nil {
		cs.storages[archetypeIndex] = []unsafe.Pointer{}
	}
	// TODO: optimize to avoid allocation
	componentValue := component.New()
	cs.storages[archetypeIndex] = append(cs.storages[archetypeIndex], componentValue)
}

// Component returns the pointer to data of the component in the archetype.
func (cs *Storage) Component(archetypeIndex ArchetypeIndex, componentIndex ComponentIndex) unsafe.Pointer {
	return cs.storages[archetypeIndex][componentIndex]
}

// SetComponent sets the pointer to data of the component in the archetype.
func (cs *Storage) SetComponent(archetypeIndex ArchetypeIndex, componentIndex ComponentIndex, component unsafe.Pointer) {
	cs.storages[archetypeIndex][componentIndex] = component
}

// MoveComponent moves the pointer to data of the component in the archetype.
func (cs *Storage) MoveComponent(source ArchetypeIndex, index ComponentIndex, dst ArchetypeIndex) {
	if len(cs.storages) <= int(dst) {
		cs.ensureCapacity(dst)
	}

	src_slice := cs.storages[source]
	dst_slice := cs.storages[dst]

	value := src_slice[index]
	src_slice[index] = src_slice[len(src_slice)-1]
	src_slice = src_slice[:len(src_slice)-1]
	cs.storages[source] = src_slice

	dst_slice = append(dst_slice, value)
	cs.storages[dst] = dst_slice
}

// SwapRemove removes the pointer to data of the component in the archetype.
func (cs *Storage) SwapRemove(archetypeIndex ArchetypeIndex, componentIndex ComponentIndex) unsafe.Pointer {
	componentValue := cs.storages[archetypeIndex][componentIndex]
	cs.storages[archetypeIndex][componentIndex] = cs.storages[archetypeIndex][len(cs.storages[archetypeIndex])-1]
	cs.storages[archetypeIndex] = cs.storages[archetypeIndex][:len(cs.storages[archetypeIndex])-1]
	return componentValue
}

// Contains returns true if the storage contains the component.
func (cs *Storage) Contains(archetypeIndex ArchetypeIndex, componentIndex ComponentIndex) bool {
	if cs.storages[archetypeIndex] == nil {
		return false
	}
	if len(cs.storages[archetypeIndex]) <= int(componentIndex) {
		return false
	}
	return cs.storages[archetypeIndex][componentIndex] != nil
}

func (cs *Storage) ensureCapacity(archetypeIndex ArchetypeIndex) {
	newStorages := make([][]unsafe.Pointer, len(cs.storages)*2)
	copy(newStorages, cs.storages)
	cs.storages = newStorages
}
