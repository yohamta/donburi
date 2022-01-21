package storage

import (
	"github.com/yohamta/donburi/internal/component"
	"github.com/yohamta/donburi/internal/entity"
)

type ArchetypeIndex int

// Archetype is a collection of entities for a specific layout of components.
// This structure allows to quickly find entities based on their components.
type Archetype struct {
	index    ArchetypeIndex
	entities []entity.Entity
	layout   *EntityLayout
}

// NewArchetype creates a new archetype.
func NewArchetype(index ArchetypeIndex, layout *EntityLayout) *Archetype {
	return &Archetype{
		index:    index,
		entities: make([]entity.Entity, 0, 256),
		layout:   layout,
	}
}

// Layout is a collection of archetypes for a specific layout of components.
func (archetype *Archetype) Layout() *EntityLayout {
	return archetype.layout
}

// Entities returns all entities in this archetype.
func (archetype *Archetype) Entities() []entity.Entity {
	return archetype.entities
}

// SwapRemove removes an entity from the archetype and returns it.
func (archetype *Archetype) SwapRemove(entity_index int) entity.Entity {
	removed := archetype.entities[entity_index]
	archetype.entities[entity_index] = archetype.entities[len(archetype.entities)-1]
	archetype.entities = archetype.entities[:len(archetype.entities)-1]
	return removed
}

// Matches returns true if the given components match the archetype.
func (archetype *Archetype) MatchesLayout(components []*component.ComponentType) bool {
	for _, componentType := range components {
		if !archetype.layout.HasComponent(componentType) {
			return false
		}
	}
	return true
}

// PushEntity adds an entity to the archetype.
func (archetype *Archetype) PushEntity(entity entity.Entity) {
	archetype.entities = append(archetype.entities, entity)
}
