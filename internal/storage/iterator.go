package storage

import "github.com/yohamta/donburi/internal/entity"

// EntityIterator is an iterator for entity lists in archetypes.
type EntityIterator struct {
	current    int
	archetypes []*Archetype
	indices    []ArchetypeIndex
}

// EntityIterator is an iterator for entities.
func NewEntityIterator(current int, archetypes []*Archetype, indices []ArchetypeIndex) EntityIterator {
	return EntityIterator{
		current:    current,
		archetypes: archetypes,
		indices:    indices,
	}
}

// HasNext returns true if there are more entity list to iterate over.
func (it *EntityIterator) HasNext() bool {
	return it.current < len(it.indices)
}

// Next returns the next entity list.
func (it *EntityIterator) Next() []entity.Entity {
	archetypeIndex := it.indices[it.current]
	it.current++
	return it.archetypes[archetypeIndex].Entities()
}
