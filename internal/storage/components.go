package storage

import (
	"github.com/yohamta/donburi/internal/component"
)

// ComponentIndex represents the index of component in a archetype.
type ComponentIndex int

// Components is a structure that stores data of components.
type Components struct {
	storages []*Storage
	// TODO: optimize to use slice instead of map for performance
	componentIndices map[ArchetypeIndex]ComponentIndex
}

// NewComponents creates a new empty structure that stores data of components.
func NewComponents() *Components {
	return &Components{
		storages:         make([]*Storage, 512), // TODO: expand as the number of component types increases
		componentIndices: make(map[ArchetypeIndex]ComponentIndex),
	}
}

// PUshComponent stores the new data of the component in the archetype.
func (cs *Components) PushComponents(components []*component.ComponentType, archetypeIndex ArchetypeIndex) ComponentIndex {
	for _, componentType := range components {
		if v := cs.storages[componentType.Id()]; v == nil {
			cs.storages[componentType.Id()] = NewStorage()
		}
		cs.storages[componentType.Id()].PushComponent(componentType, archetypeIndex)
	}
	if _, ok := cs.componentIndices[archetypeIndex]; !ok {
		cs.componentIndices[archetypeIndex] = 0
	} else {
		cs.componentIndices[archetypeIndex]++
	}
	return cs.componentIndices[archetypeIndex]
}

// MoveComponent moves the pointer to data of the component in the archetype.
func (cs *Components) MoveComponent(c *component.ComponentType, src ArchetypeIndex, idx ComponentIndex, dst ArchetypeIndex) {
	storage := cs.Storage(c)
	storage.MoveComponent(src, idx, dst)
	cs.componentIndices[src]--
}

// Storage returns the pointer to data of the component in the archetype.
func (cs *Components) Storage(c *component.ComponentType) *Storage {
	if storage := cs.storages[c.Id()]; storage != nil {
		return storage
	}
	cs.storages[c.Id()] = NewStorage()
	return cs.storages[c.Id()]
}

// Remove removes the component from the storage.
func (cs *Components) Remove(a *Archetype, ci ComponentIndex) {
	for _, ct := range a.Layout().Components() {
		cs.remove(ct, a.index, ci)
	}
	cs.componentIndices[a.index]--
}

func (cs *Components) remove(ct *component.ComponentType,
	ai ArchetypeIndex, ci ComponentIndex) {
	storage := cs.Storage(ct)
	storage.SwapRemove(ai, ci)
}
