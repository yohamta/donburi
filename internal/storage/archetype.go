package storage

import (
	"github.com/yohamta/donburi/component"
	"sync"
)

type ArchetypeIndex int

// Archetype is a collection of entities for a specific layout of components.
// This structure allows to quickly find entities based on their components.
type Archetype struct {
	index    ArchetypeIndex
	entities []Entity
	layout   *Layout
	lock     sync.Mutex
}

func (archetype *Archetype) Lock() {
	archetype.lock.Lock()
}

func (archetype *Archetype) Unlock() {
	archetype.lock.Unlock()
}

// NewArchetype creates a new archetype.
func NewArchetype(index ArchetypeIndex, layout *Layout) *Archetype {
	return &Archetype{
		index:    index,
		entities: make([]Entity, 0, 256),
		layout:   layout,
	}
}

// Layout is a collection of archetypes for a specific layout of components.
func (archetype *Archetype) Layout() *Layout {
	return archetype.layout
}

// Entities returns all entities in this archetype.
func (archetype *Archetype) Entities() []Entity {
	return archetype.entities
}

// ComponentTypes returns a slice of all component types in the archetype.
func (archetype *Archetype) ComponentTypes() []component.IComponentType {
	return archetype.Layout().componentTypes
}

// SwapRemove removes an entity from the archetype and returns it.
func (archetype *Archetype) SwapRemove(entityIndex int) Entity {
	removed := archetype.entities[entityIndex]
	archetype.entities[entityIndex] = archetype.entities[len(archetype.entities)-1]
	archetype.entities = archetype.entities[:len(archetype.entities)-1]
	return removed
}

// LayoutMatches returns true if the given layout matches this archetype.
func (archetype *Archetype) LayoutMatches(components []component.IComponentType) bool {
	if len(archetype.layout.Components()) != len(components) {
		return false
	}
	for _, componentType := range components {
		if !archetype.layout.HasComponent(componentType) {
			return false
		}
	}
	return true
}

// PushEntity adds an entity to the archetype.
func (archetype *Archetype) PushEntity(entity Entity) {
	archetype.entities = append(archetype.entities, entity)
}

// Count returns the number of entities in the archetype.
func (archetype *Archetype) Count() int {
	return len(archetype.entities)
}
