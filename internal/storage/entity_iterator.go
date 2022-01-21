package storage

import "github.com/yohamta/donburi/internal/entity"

// EntityListIterator is an iterator for entity lists in archetypes.
type EntityListIterator struct {
	current    int
	archetypes []*Archetype
	indices    []ArchetypeIndex
}

// EntityIterator is an iterator for entities.
func NewEntityListIterator(current int, archetypes []*Archetype, indices []ArchetypeIndex) EntityListIterator {
	return EntityListIterator{
		current:    current,
		archetypes: archetypes,
		indices:    indices,
	}
}

// HasNext returns true if there are more entity list to iterate over.
func (it *EntityListIterator) HasNext() bool {
	return it.current < len(it.indices)
}

// Next returns the next entity list.
func (it *EntityListIterator) Next() []entity.Entity {
	archetypeIndex := it.indices[it.current]
	it.current++
	return it.archetypes[archetypeIndex].Entities()
}
