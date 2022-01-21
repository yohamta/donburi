package storage

import (
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/donburi/internal/component"
)

type archetypeIterator struct {
	current int
	values  []ArchetypeIndex
}

func (it *archetypeIterator) HasNext() bool {
	return it.current < len(it.values)
}

func (it *archetypeIterator) Next() ArchetypeIndex {
	val := it.values[it.current]
	it.current++
	return val
}

// SearchIndex is a structure that indexes archetypes by their component types.
type SearchIndex struct {
	component_layouts [][]*component.ComponentType
	iterator          *archetypeIterator
}

// NewSearchIndex creates a new search index.
func NewSearchIndex() *SearchIndex {
	return &SearchIndex{
		component_layouts: [][]*component.ComponentType{},
		iterator: &archetypeIterator{
			current: 0,
		},
	}
}

// Push adds an archetype to the search index.
func (si *SearchIndex) Push(archetypeLayout *EntityLayout) {
	si.component_layouts = append(si.component_layouts, archetypeLayout.Components())
}

// SearchFrom searches for archetypes that match the given filter from the given index.
func (si *SearchIndex) SearchFrom(f filter.LayoutFilter, start int) *archetypeIterator {
	si.iterator.current = 0
	si.iterator.values = []ArchetypeIndex{}
	for i := start; i < len(si.component_layouts); i++ {
		if f.MatchesLayout(si.component_layouts[i]) {
			si.iterator.values = append(si.iterator.values, ArchetypeIndex(i))
		}
	}
	return si.iterator
}

// Search searches for archetypes that match the given filter.
func (si *SearchIndex) Search(filter filter.LayoutFilter) *archetypeIterator {
	return si.SearchFrom(filter, 0)
}
