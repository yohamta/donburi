package storage

import (
	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/filter"
)

type ArchetypeIterator struct {
	current int
	values  []ArchetypeIndex
}

func (it *ArchetypeIterator) HasNext() bool {
	return it.current < len(it.values)
}

func (it *ArchetypeIterator) Next() ArchetypeIndex {
	val := it.values[it.current]
	it.current++
	return val
}

// Index is a structure that indexes archetypes by their component types.
type Index struct {
	layouts  [][]component.IComponentType
	iterator *ArchetypeIterator
}

// NewIndex creates a new search index.
func NewIndex() *Index {
	return &Index{
		layouts: [][]component.IComponentType{},
		iterator: &ArchetypeIterator{
			current: 0,
		},
	}
}

// Push adds an archetype to the search index.
func (idx *Index) Push(layout *Layout) {
	idx.layouts = append(idx.layouts, layout.Components())
}

// SearchFrom searches for archetypes that match the given filter from the given index.
func (idx *Index) SearchFrom(f filter.LayoutFilter, start int) *ArchetypeIterator {
	idx.iterator.current = 0
	idx.iterator.values = []ArchetypeIndex{}
	for i := start; i < len(idx.layouts); i++ {
		if f.MatchesLayout(idx.layouts[i]) {
			idx.iterator.values = append(idx.iterator.values, ArchetypeIndex(i))
		}
	}
	return idx.iterator
}

// Search searches for archetypes that match the given filter.
func (idx *Index) Search(filter filter.LayoutFilter) *ArchetypeIterator {
	return idx.SearchFrom(filter, 0)
}
