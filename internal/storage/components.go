package storage

import "github.com/yohamta/donburi/component"

// ComponentIndex represents the index of component in an archetype.
type ComponentIndex int

// Components is a structure that stores data of components.
type Components struct {
	storages         []*Storage
	componentIndices map[ArchetypeIndex]ComponentIndex
}

// NewComponents creates a new empty structure that stores data of components.
func NewComponents() *Components {
	return &Components{
		storages:         make([]*Storage, 0, 512),
		componentIndices: make(map[ArchetypeIndex]ComponentIndex),
	}
}

// PushComponents stores the new data of the component in the archetype.
func (cs *Components) PushComponents(components []component.IComponentType, archetypeIndex ArchetypeIndex) ComponentIndex {
	for _, componentType := range components {
		id := int(componentType.Id())
		if id >= len(cs.storages) {
			cs.storages = append(cs.storages, make([]*Storage, id-len(cs.storages)+1)...)
		}
		if cs.storages[id] == nil {
			cs.storages[id] = NewStorage()
		}
		cs.storages[id].PushComponent(componentType, archetypeIndex)
	}
	if _, ok := cs.componentIndices[archetypeIndex]; !ok {
		cs.componentIndices[archetypeIndex] = 0
	} else {
		cs.componentIndices[archetypeIndex]++
	}
	return cs.componentIndices[archetypeIndex]
}

// MoveComponent moves the pointer to data of the component in the archetype.
func (cs *Components) Move(src ArchetypeIndex, dst ArchetypeIndex) {
	cs.componentIndices[src]--
	cs.componentIndices[dst]++
}

// Storage returns the pointer to data of the component in the archetype.
func (cs *Components) Storage(c component.IComponentType) *Storage {
	id := int(c.Id())
	if id < len(cs.storages) && cs.storages[id] != nil {
		return cs.storages[id]
	}
	if id >= len(cs.storages) {
		cs.storages = append(cs.storages, make([]*Storage, id-len(cs.storages)+1)...)
	}
	cs.storages[id] = NewStorage()
	return cs.storages[id]
}

// Remove removes the component from the storage.
func (cs *Components) Remove(a *Archetype, ci ComponentIndex) {
	for _, ct := range a.Layout().Components() {
		cs.remove(ct, a.index, ci)
	}
	cs.componentIndices[a.index]--
}

func (cs *Components) remove(ct component.IComponentType, ai ArchetypeIndex, ci ComponentIndex) {
	storage := cs.Storage(ct)
	storage.SwapRemove(ai, ci)
}
